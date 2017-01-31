// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017 Peter Angyal
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"github.com/angyalp/homecontrol/messaging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const weather_poll_sec int = 60

func send_sensor_value(bus *messaging.Bus, sensor_id string, sensor_value interface{}) {
	err := bus.SendEvent(messaging.SensorInterface, messaging.SensorValue, sensor_id, sensor_value)
	if err != nil {
		fmt.Printf("Unable to send sensor value! Sensor/Value: %s/%v, Error: %v\n", sensor_id, sensor_value, err)
	}
}

func send_sensor_event(bus *messaging.Bus, sensor_id string) {
	err := bus.SendEvent(messaging.SensorInterface, messaging.SensorEvent, sensor_id)
	if err != nil {
		fmt.Printf("Unable to send sensor event! Event: %s, Error: %v\n", sensor_id, err)
	}
}

func update_weather(bus *messaging.Bus) {
	fmt.Println("Updating weather...")
	send_sensor_value(bus, "web1.temp", uint32(10))
	send_sensor_value(bus, "web1.hum", uint32(50))
	send_sensor_event(bus, "web1.lightning")
}

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	bus, err := messaging.RegisterApp("com.github.homecontrol.webweather")
	if err != nil {
		fmt.Println("Unable to create message bus! Error:", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Second * time.Duration(weather_poll_sec))
	go func() {
		for {
			update_weather(bus)
			<-ticker.C
		}
	}()

	select {
	case sig := <-ch:
		fmt.Println("Received signal:", sig)
	}
}
