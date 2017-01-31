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
)

const bus_event_channel_length int = 1000

func handle_event(event *messaging.Event) {
	fmt.Println(event)
	if event.GetInterfaceName() == messaging.SensorInterface {
		switch event.GetEventName() {
		case messaging.SensorValue:
			{
				// TODO check Body length. Can DBus enforce its format?
				fmt.Printf("Sensor value: %v:%v\n", event.Body[0], event.Body[1])
			}
		case messaging.SensorEvent:
			{
				// TODO check Body length. Can DBus enforce its format?
				fmt.Printf("Sensor event: %v\n", event.Body[0])
			}
		}
	}
}

func main() {
	signal_ch := make(chan os.Signal)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	bus, err := messaging.RegisterApp("com.github.homecontrol.dw")
	if err != nil {
		fmt.Println("Unable to create message bus! Error:", err)
		os.Exit(1)
	}

	err = bus.RegisterForEvent(messaging.SensorInterface, "")
	if err != nil {
		fmt.Println("Unable to register for sensor events! Error:", err)
		os.Exit(1)
	}

	event_channel := bus.GetEventChannel(bus_event_channel_length)

	go func() {
		for event := range event_channel {
			handle_event((*messaging.Event)(event))
		}
	}()

	select {
	case os_signal := <-signal_ch:
		fmt.Println("Received signal:", os_signal)
	}
}
