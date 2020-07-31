package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aler9/goroslib/msgs/geometry_msgs"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"os/exec"
)

const projectID = "roswheel"
const topicName = "ctrl"

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	c := client.Topic(topicName)

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var b []byte = make([]byte, 1)
	msg := &geometry_msgs.Twist{
		Linear:  geometry_msgs.Vector3{},
		Angular: geometry_msgs.Vector3{},
	}
	fmt.Println("initialized")
	for {
		os.Stdin.Read(b)
		switch string(b) {
		case "w":
			msg.Linear = geometry_msgs.Vector3{X: 10, Y: 0}
		case "a":
			msg.Linear = geometry_msgs.Vector3{X: 0, Y: -10}
		case "s":
			msg.Linear = geometry_msgs.Vector3{X: -10, Y: 0}
		case "d":
			msg.Linear = geometry_msgs.Vector3{X: 0, Y: 10}
		}

		b, err := msgpack.Marshal(msg)
		if err != nil {
			panic(err)
		}
		c.Publish(ctx, &pubsub.Message{
			Data: b,
		})
	}
}
