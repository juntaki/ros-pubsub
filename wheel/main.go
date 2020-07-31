package main

import (
	"fmt"
	"github.com/aler9/goroslib/msgs/geometry_msgs"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"time"

	"github.com/aler9/goroslib"
)

var wheelBase = 11.8
var tred = 13.3

func main() {
	a := raspi.NewAdaptor()
	d := i2c.NewAdafruitMotorHatDriver(
		a,
	)
	d.SetMotorHatAddress(0x6f)
	d.Start()

	ti := time.NewTimer(1 * time.Second)
	go func() {
		for {
			<-ti.C
			d.RunDCMotor(0, i2c.AdafruitRelease)
			d.RunDCMotor(1, i2c.AdafruitRelease)
			d.RunDCMotor(2, i2c.AdafruitRelease)
			d.RunDCMotor(3, i2c.AdafruitRelease)
		}
	}()

	onMessage := func(msg *geometry_msgs.Twist) {
		vx := msg.Linear.X
		vy := msg.Linear.Y
		om := msg.Angular.Z
		a := (wheelBase/2 + tred/2) / 25
		velof := make([]float64, 4)

		velof[0] = vx - vy - a*om
		velof[1] = vx + vy - a*om
		velof[2] = vx - vy + a*om
		velof[3] = vx + vy + a*om

		fmt.Println(velof)

		for i, vf := range velof {
			v := int32(math.Abs(vf) * 50)
			if v > 255 {
				v = 255
			}
			d.SetDCMotorSpeed(i, v)
			if vf > 0 {
				d.RunDCMotor(i, i2c.AdafruitForward)
			} else {
				d.RunDCMotor(i, i2c.AdafruitBackward)
			}
		}

		if !ti.Stop() {
			select {
			case <-ti.C:
			default:
			}
		}
		ti.Reset(500 * time.Millisecond)
	}

	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:       "/robo",
		MasterHost: "127.0.0.1",
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     n,
		Topic:    "/robo/cmd_vel",
		Callback: onMessage,
		Protocol: 0,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// freeze main loop
	select {}
}
