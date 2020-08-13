package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"time"

	_ "github.com/anrid/roboviewer/docs" // Swag way of doing things.
	docs "github.com/anrid/roboviewer/docs"
	"github.com/anrid/roboviewer/robo/config"
	"github.com/anrid/roboviewer/robo/controller"
	"github.com/anrid/roboviewer/robo/dg"
	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/anrid/roboviewer/robo/pkg/mqtt"
	"github.com/anrid/roboviewer/robo/pkg/msgdel"
	"github.com/anrid/roboviewer/robo/service"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server is an API server.
type Server struct{}

// Run ...
func Run() {
	ctx := context.Background()

	log.Printf("setup and run API server")

	var err error
	c := config.GetConfig()

	// Dump server config.
	entity.Dump(c)

	// Connect to Dgraph.
	conn, disconnect := dg.Connect(c.DgraphURL)
	defer disconnect()

	if c.DropAll {
		// Drop and recreate database schema and also seed
		// some test data.
		println("Dropping and recreating database schema ..")
		dg.DropAll(ctx, conn)
		dg.CreateSchema(ctx, conn)
		dg.CreateSimpleTestData(ctx, conn)
		println("Done.")
		os.Exit(0)
	}

	if c.Migrate {
		println("Applying database migrations ..")
		dg.CreateSchema(ctx, conn)
	}

	// Setup repositories.
	repos := struct {
		Robot entity.RobotRepository
		Area  entity.AreaRepository
	}{
		Robot: dg.NewRobotRepository(conn),
		Area:  dg.NewAreaRepository(conn),
	}

	// Setup services.
	svcs := struct {
		Robot entity.RobotService
		Area  entity.AreaService
	}{
		Robot: service.NewRobotService(repos.Robot),
		Area:  service.NewAreaService(repos.Area),
	}

	// New HTTP server.
	serv := httpserver.NewServer()

	// Setup controllers.
	controller.NewRobotController(svcs.Robot).SetupRoutes(serv.Echo)
	controller.NewAreaController(svcs.Area).SetupRoutes(serv.Echo)

	// Wire up our message delegator to MQTT broker to handle
	// incoming MQTT messages from robots.
	broker := mqtt.NewClient(c.MQTTBrokerURL)

	delegator := msgdel.NewMessageDelegator(svcs.Robot)

	broker.Subscribe(c.TopicRobotSessionStart, delegator.HandleStartSession)
	broker.Subscribe(c.TopicRobotSessionUpdate, delegator.HandleUpdateSession)
	broker.Subscribe(c.TopicRobotSessionEnd, delegator.HandleEndSession)

	// Setup Swagger documentation.
	docs.SwaggerInfo.Host = c.Host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	serv.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		func(cfg *echoSwagger.Config) {
			cfg.URL = c.APIURL + "/swagger/doc.json"
		}),
	)
	log.Printf("swagger url: %s", c.APIURL+"/swagger/index.html")

	// Dump routes to stdout.
	{
		var routes []*echo.Route
		for _, r := range serv.Echo.Routes() {
			if r.Path != "" && r.Path != "/*" {
				routes = append(routes, r)
			}
		}
		sort.SliceStable(routes, func(i, j int) bool { return routes[i].Path < routes[j].Path })
		for i, r := range routes {
			fmt.Printf("route %03d: [%-6s]  %s\n", i+1, r.Method, r.Path)
		}
	}

	// Start server.
	go func() {
		log.Printf("starting API server at %s", c.APIURL)
		err = serv.Echo.Start(c.Host)
		if err != nil {
			log.Fatalf("could not start API server: %s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the
	// server with a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := serv.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown API server gracefully: %s", err.Error())
	}
	log.Printf("oh, bye!")

	return
}
