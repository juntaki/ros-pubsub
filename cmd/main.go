package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/msgs/geometry_msgs"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/api/option"
)

const projectID = "roswheel"
const subscriptionName = "robot"

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("cred.json"))
	if err != nil {
		panic(err)
	}

	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:       "/robo_cmd",
		MasterHost: "127.0.0.1",
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	pub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  n,
		Topic: "/robo/cmd_vel",
		Msg:   &geometry_msgs.Twist{},
	})
	if err != nil {
		panic(err)
	}
	defer pub.Close()

	s := client.Subscription(subscriptionName)

	fmt.Println("initialized")
	_ = s.Receive(ctx, func(c context.Context, m *pubsub.Message) {
		var msg geometry_msgs.Twist
		err = msgpack.Unmarshal(m.Data, &msg)
		if err != nil {
			panic(err)
		}
		pub.Write(&msg)
	})
}
