# Robo Viewer

## Features

The backend demo is built as follows;

- [Dgraph](https://dgraph.io/) as database.
- [Echo](https://github.com/labstack/echo) as HTTP server.
- [Mosquitto](http://mosquitto.org/) for our MQTT broker.
- The project is structured around the ideas behind [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

## Run Demo

Start the server as follows:

```bash
# Setup docker volumes:
./setup-volumes.sh

# Start database and MQTT broker:
docker-compose up -d

# Setup database schema and seed with test data:
go run cmd/server/main.go -drop-all

# Run load test to generate more data:
go run cmd/loadtest/main.go -concur 2

# Start backend API server at http://localhost:3000:
go run cmd/server/main.go
```

Then in a new terminal window:

```bash
# Root / health check:
curl http://localhost:3000
# OUTPUT: {"message":"All is well in the world!","ok":true,"timestamp":1581828959234514000}

# List robots:
curl http://localhost:3000/v1/robots
# OUTPUT: {"ok":true,"robots":[{"name":"Test - Johnny 5","session":[...

# Show robot history:
curl http://localhost:3000/v1/robots/0x64/history
# OUTPUT: {"ok":true,"robot":{"name":"Test - Johnny 5","session":[...
```

## Config

```bash
go run cmd/server/main.go -help
# Prints:
#
#  -drop-all
#    	drop all tables and recreate schema
#  -migrate
#    	migrate schema changes
#  -mqtt-broker-url string
#    	set MQTT broker URL, e.g tcp://localhost:1883 (default "tcp://localhost:1883")
#  -topic-end string
#    	set MQTT topic for cleaning session end (default "/robot/session/end")
#  -topic-start string
#    	set MQTT topic for cleaning session start (default "/robot/session/start")
#  -topic-update string
#    	set MQTT topic for robot session update (default "/robot/session/update")
#
```

## Swagger Documentation

Run server then goto:

http://localhost:3000/swagger/index.html

## Run Load Test

```bash
# Ensure database and MQTT broker are running (see Run Demo) then:
go run cmd/loadtest/main.go -concur 50 # Runs 50 concurrent robot cleaning sessions.
```

## Run Developer Tests

```bash
# Ensure database and MQTT broker are running (see Run Demo), then:
./test.sh
```

## Migrate database

Update schema in `robo/dg/schema.go` and run:

```bash
go run cmd/server/main.go -migrate
```

## Background

The initial spec discussion can be found here:
https://docs.google.com/document/d/1QXxDKCEnrcUDR8a2sk4pq5_i5HYfpIwkgZccuqGnqSo/edit

## Spec

### 1. The platform must enable the users to:

#### 1.1. Fetch the current reported position of the robot.

- Positional data will be a simple x,y coordinate indicating
  distance from top left corner of the grid.

#### 1.2. Fetch the current state of completion.

- Completion can be a simple percentage.

- Area data is coverted into a grid and the robot may have to visit
  each grid square multiple times before that square can be considered
  clean.

### 2. The platform must enable the robots to:

#### 2.1. To send the positional and status data to the backend.

### 3. The platform might enable the users to:

#### 3.1. Fetch historical data

- 3.1.1. With the paths taken.
- 3.1.2. With completion time.
- 3.1.3. Approximated coverage.

### TODO

- [x] _TAS-1_: Add Mosquitto MQTT server and basic pub/sub test.
- [x] _TAS-2_: Add types for robots, robot cleaning sessions, positions and cleaning area.
- [x] _TAS-3_: Refactor area type, add grid, add support for multiple visits to the same grid square.
- [x] _TAS-4_: Init database schema and generate some test data based on the types created in _TAS-2_.
- [x] _TAS-5_: Create robot repository and try reading / writing a sizeable object graph.
- [x] _TAS-6_: Refactor types be more suitable for Dgraph and improve database tests.
- [x] _TAS-7_: Add simple http server to receive REST API calls.
- [x] _TAS-8_: Add service and controller layers.
- [x] _TAS-9_: Add support for starting new cleaning sessions and updating an ongoing cleaning session.
- [x] _TAS-10_: Create load test that pub/subs on all MQTT topics that robots will use.
- [x] _TAS-11_: Add missing REST endpoints to display session history for a particular robot.
- [x] _TAS-12_: Subscribe to MQTT broker and add missing documentation.
