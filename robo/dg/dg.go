package dg

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// CancelFunc defines the function to call when we want
// to disconnect from Dgraph.
type CancelFunc func()

// Repository contains low-level database operations.
type Repository struct {
	c *dgo.Dgraph
}

// Save the given object.
func (r *Repository) Save(ctx context.Context, object interface{}) (map[string]string, error) {
	return Store(ctx, r.c, object)
}

// Run is a basic sanity check.
func Run(dgraphURL string, dropAll bool) {
	ctx := context.Background()
	c, cancel := Connect(dgraphURL)
	defer cancel()

	if dropAll {
		DropAll(ctx, c)
		CreateSchema(ctx, c)
	}

	res1 := CreateSimpleTestData(ctx, c)
	spew.Dump(res1)

	rr := NewRobotRepository(c)
	ar := NewAreaRepository(c)

	res2, err := rr.List(ctx, entity.ListRobotsArgs{Name: "test"})
	if err != nil {
		log.Panicf("could not list robots: %s", err.Error())
	}
	for _, robot := range res2.Robots {
		println("robot:", robot.UID, robot.Name)
	}

	res3, err := ar.List(ctx)
	if err != nil {
		log.Panicf("could not list areas: %s", err.Error())
	}
	for _, area := range res3.Areas {
		println("area:", area.UID, area.Name)
	}

	println("It’s a Done Deal!")
}

// Connect to Dgraph server.
func Connect(dgraphURL string) (*dgo.Dgraph, CancelFunc) {
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to Dgraph using gRPC")
	}

	dc := api.NewDgraphClient(conn)
	c := dgo.NewDgraphClient(dc)

	// Perform login call. If the Dgraph cluster does not have ACL and
	// enterprise features enabled, this call should be skipped.
	// ctx := context.Background()
	// for {
	// 	// Keep retrying until we succeed or receive a non-retriable error.
	// 	err = dg.Login(ctx, "groot", "password")
	// 	if err == nil || !strings.Contains(err.Error(), "Please retry") {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }
	// if err != nil {
	// 	log.Fatalf("While trying to login %v", err.Error())
	// }

	return c, func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("error while closing connection: %s", err.Error())
		}
	}
}

// DropAll drops all schema and type definitions.
func DropAll(ctx context.Context, c *dgo.Dgraph) {
	op := api.Operation{DropAll: true}
	if err := c.Alter(ctx, &op); err != nil {
		log.Fatalf("could not perform drop all: %s", err.Error())
	}
}

// Store persists an object graph.
func Store(ctx context.Context, c *dgo.Dgraph, d interface{}) (map[string]string, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("could not marshal query: %s", err.Error())
	}
	mu.SetJson = b

	// Execute.
	res, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	log.Printf("set node in %d ms", res.GetLatency().ProcessingNs/uint64(time.Millisecond))

	// New uids for nodes which were created are returned in the
	// response.Uids map.
	return res.Uids, nil
}

// Edge creates an edge between two objects.
func Edge(ctx context.Context, c *dgo.Dgraph, data ...string) error {
	mu := &api.Mutation{
		CommitNow: true,
	}

	var quads bytes.Buffer
	var edges int
	for i := 0; i < len(data); i += 3 {
		if i+2 < len(data) {
			quads.WriteString(fmt.Sprintf("<%s> <%s> <%s> .\n", data[i], data[i+1], data[i+2]))
			edges++
		}
	}
	// println("edges:")
	// println(string(quads.Bytes()))
	mu.SetNquads = quads.Bytes()

	// Execute.
	res, err := c.NewTxn().Mutate(ctx, mu)

	log.Printf("set %d edges in %d ms", edges, res.GetLatency().ProcessingNs/uint64(time.Millisecond))
	return err
}
