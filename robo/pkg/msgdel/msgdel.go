package msgdel

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anrid/roboviewer/robo/entity"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ entity.MessageDelegator = &MessageDelegator{}

// MessageDelegator implements entity.MessageDelegator.
// It takes an MQTT message, converts it, and passes it
// on to the RobotService.
type MessageDelegator struct {
	svc entity.RobotService
}

// NewMessageDelegator creates a new MessageDelegator instance.
func NewMessageDelegator(svc entity.RobotService) *MessageDelegator {
	return &MessageDelegator{svc}
}

// HandleStartSession handles incoming start cleaning session messages
// from robots.
func (md *MessageDelegator) HandleStartSession(c mqtt.Client, m mqtt.Message) {
	msg := string(m.Payload())

	parts := strings.SplitN(msg, "/", 5)
	if len(parts) != 5 {
		log.Printf("invalid message '%s', should contain 'robotID/areaID/robotX/robotY/unixTimestamp'", msg)
		return
	}

	robotID := parts[0]
	areaID := parts[1]

	x, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("invalid x coordinate in message: '%s'", msg)
	}
	y, err := strconv.Atoi(parts[3])
	if err != nil {
		log.Printf("invalid y coordinate in message: '%s'", msg)
	}
	ts, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		log.Printf("invalid timestamp coordinate in message: '%s'", msg)
	}

	startedAt := time.Unix(ts, 0)

	sess, err := md.svc.StartSession(context.Background(), entity.StartSessionArgs{
		RobotID:   robotID,
		AreaID:    areaID,
		RobotX:    x,
		RobotY:    y,
		StartedAt: startedAt,
	})
	if err != nil {
		log.Printf("could not start session: %s", err.Error())
		return
	}

	log.Printf("started cleaning session: %s %s", sess.UID, sess.Name)
}

// HandleUpdateSession handles incoming update cleaning session messages
// from robots, e.g. when a robot moves.
func (md *MessageDelegator) HandleUpdateSession(c mqtt.Client, m mqtt.Message) {
	msg := string(m.Payload())

	parts := strings.SplitN(msg, "/", 4)
	if len(parts) != 4 {
		log.Printf("invalid message '%s', should contain 'robotID/robotX/robotY/unixTimestamp'", msg)
		return
	}

	robotID := parts[0]

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("invalid x coordinate in message: '%s'", msg)
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("invalid y coordinate in message: '%s'", msg)
	}
	ts, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Printf("invalid timestamp coordinate in message: '%s'", msg)
	}

	startedAt := time.Unix(ts, 0)

	sess, err := md.svc.UpdateSession(context.Background(), entity.UpdateSessionArgs{
		RobotID:    robotID,
		RobotX:     x,
		RobotY:     y,
		ReportedAt: startedAt,
	})
	if err != nil {
		log.Printf("could not update session: %s", err.Error())
		return
	}

	log.Printf("updated cleaning session: %s %s", sess.UID, sess.Name)
}

// HandleEndSession handles incoming end cleaning session messages
// from robots.
func (md *MessageDelegator) HandleEndSession(c mqtt.Client, m mqtt.Message) {
	msg := string(m.Payload())

	parts := strings.SplitN(msg, "/", 4)
	if len(parts) != 4 {
		log.Printf("invalid message '%s', should contain 'robotID/robotX/robotY/unixTimestamp'", msg)
		return
	}

	robotID := parts[0]

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("invalid x coordinate in message: '%s'", msg)
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("invalid y coordinate in message: '%s'", msg)
	}
	ts, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Printf("invalid timestamp coordinate in message: '%s'", msg)
	}

	startedAt := time.Unix(ts, 0)

	sess, err := md.svc.EndSession(context.Background(), entity.UpdateSessionArgs{
		RobotID:    robotID,
		RobotX:     x,
		RobotY:     y,
		ReportedAt: startedAt,
	})
	if err != nil {
		log.Printf("could not end session: %s", err.Error())
		return
	}

	log.Printf("ended cleaning session: %s %s", sess.UID, sess.Name)
}
