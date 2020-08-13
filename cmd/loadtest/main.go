package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/anrid/roboviewer/robo/config"
	"github.com/anrid/roboviewer/robo/dg"
	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/mqtt"
	"github.com/anrid/roboviewer/robo/pkg/msgdel"
	"github.com/anrid/roboviewer/robo/service"
)

func main() {
	concur := flag.Int("concur", 4, "number of concurrent cleaning sessions to run")
	flag.Parse()

	ctx := context.Background()

	c := config.GetConfig()

	conn, _ := dg.Connect(c.DgraphURL)

	_ = dg.CreateSimpleTestData(context.Background(), conn)

	robotRepo := dg.NewRobotRepository(conn)
	areaRepo := dg.NewAreaRepository(conn)

	robotSvc := service.NewRobotService(robotRepo)
	areaSvc := service.NewAreaService(areaRepo)

	del := msgdel.NewMessageDelegator(robotSvc)

	broker := mqtt.NewClient(c.MQTTBrokerURL)
	broker.Subscribe(c.TopicRobotSessionStart, del.HandleStartSession)
	broker.Subscribe(c.TopicRobotSessionUpdate, del.HandleUpdateSession)
	broker.Subscribe(c.TopicRobotSessionEnd, del.HandleEndSession)

	robots, err := robotSvc.List(ctx, "", "")
	if err != nil {
		panic("err")
	}

	areas, err := areaSvc.List(ctx)
	if err != nil {
		panic("err")
	}

	println("robots:", len(robots))
	println("areas:", len(areas))

	// Simulate a number of concurrent cleaning
	// sessions.
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	for i := 0; i < *concur; i++ {
		sessionID := fmt.Sprintf("sess%03d", i+1)

		// Select a robot at random.
		robot := robots[rand.Intn(len(robots))]

		// Select an area at random.
		area := areas[rand.Intn(len(areas))]

		// Give us a random number of robot movements.
		moves := rand.Intn(50) + 2 // Min 2 moves, max 51.

		wg.Add(1)
		go newCleaningSession(sessionID, &wg, c, robot, area, moves)

		time.Sleep(500 * time.Millisecond) // Wait for 500ms before firing off the next cleaning session.
	}

	wg.Wait()

	println("waiting for subscribers to finish ...")
	time.Sleep(2000 * time.Millisecond)

	println("It’s a Done Deal.")
}

func newCleaningSession(id string, wg *sync.WaitGroup, c config.Config, r *entity.Robot, a *entity.Area, moves int) {
	defer wg.Done()

	rc := mqtt.NewClient(c.MQTTBrokerURL)

	println(id, "robot", r.UID, "start")
	rc.Publish(c.TopicRobotSessionStart, startSessionMessage(r.UID, a.UID, 0, 0, time.Now().Unix()))
	time.Sleep(1 * time.Second)

	path := NewRobotSnakePath(200, a.SizeX, a.SizeY, r.Size)
	var x int
	var y int

	for move := 1; move <= moves; move++ {
		x, y = path.NextPosition()

		println(id, "robot", r.UID, "move", move, "update", x, y)
		rc.Publish(c.TopicRobotSessionUpdate, updateSessionMessage(r.UID, x, y, time.Now().Unix()))
		time.Sleep(1 * time.Second)
	}

	println(id, "robot", r.UID, "end", x, y)
	rc.Publish(c.TopicRobotSessionEnd, endSessionMessage(r.UID, x, y, time.Now().Unix()))
}

func startSessionMessage(robotID, areaID string, x, y int, unixTimestamp int64) string {
	return fmt.Sprintf("%s/%s/%d/%d/%d", robotID, areaID, x, y, unixTimestamp)
}

func updateSessionMessage(robotID string, x, y int, unixTimestamp int64) string {
	return fmt.Sprintf("%s/%d/%d/%d", robotID, x, y, unixTimestamp)
}

func endSessionMessage(robotID string, x, y int, unixTimestamp int64) string {
	return fmt.Sprintf("%s/%d/%d/%d", robotID, x, y, unixTimestamp)
}

// RobotSnakePath creates a snake pattern, i.e. left to right, 1 down,
// right to left, 1 down, then reverse back up again.
// This path is created within the area given, for a robot with a
// certain diameter.
type RobotSnakePath struct {
	speed       int
	robotSize   int
	robotCenter int
	x           int
	y           int
	goRight     bool
	goDown      bool
	sizeX       int
	sizeY       int
}

// NewRobotSnakePath creates a new RobotSnakePath.
func NewRobotSnakePath(speed, sizeX, sizeY, robotSize int) *RobotSnakePath {
	return &RobotSnakePath{
		speed:       speed,
		sizeX:       sizeX,
		sizeY:       sizeY,
		robotSize:   robotSize,
		robotCenter: robotSize / 2,
		goRight:     true,
		goDown:      true,
	}
}

// NextPosition returns the next position along the robot's path.
func (p *RobotSnakePath) NextPosition() (int, int) {
	if p.goRight {
		p.x += p.speed
		if p.x > (p.sizeX - p.robotCenter) {
			// Are we past the right wall?
			p.x = p.sizeX - p.robotCenter
			p.goRight = false

			// "Teleport" one row down or up.
			if p.goDown {
				p.y += p.robotSize
				if p.y > (p.sizeY - p.robotCenter) {
					p.y = p.sizeY - p.robotCenter
					p.goDown = !p.goDown
				}
			} else {
				p.y -= p.robotSize
				if p.y < p.robotCenter {
					p.y = p.robotCenter
					p.goDown = !p.goDown
				}
			}
		}
	} else {
		// Subtract from position instead.
		p.x -= p.speed
		if p.x < p.robotCenter {
			// Are we past the left wall?
			p.x = p.robotCenter
			p.goRight = true

			// "Teleport" one row down or up.
			if p.goDown {
				p.y += p.robotSize
				if p.y > (p.sizeY - p.robotCenter) {
					p.y = p.sizeY - p.robotCenter
					p.goDown = !p.goDown
				}
			} else {
				p.y -= p.robotSize
				if p.y < p.robotCenter {
					p.y = p.robotCenter
					p.goDown = !p.goDown
				}
			}
		}
	}
	return p.x, p.y
}
