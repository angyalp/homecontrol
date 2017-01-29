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

package messaging

import (
	"fmt"
	"github.com/godbus/dbus"
	"strings"
)

const path_root dbus.ObjectPath = "/com/github/homecontrol"

const (
	SensorInterface   string = "com.github.homecontrol.sensor"
	ValueChangedEvent string = "value_changed"
)

type Bus struct {
	dbus_conn *dbus.Conn
}

type Event dbus.Signal

func (event *Event) GetInterfaceName() string {
	if pos := strings.LastIndex(event.Name, "."); pos == -1 {
		return event.Name
	} else {
		return event.Name[:pos]
	}
}

func (event *Event) GetEventName() string {
	if pos := strings.LastIndex(event.Name, "."); pos == -1 {
		return ""
	} else {
		if pos+1 >= len(event.Name) {
			return ""
		} else {
			return event.Name[pos+1:]
		}
	}
}

func RegisterApp(application_name string) (*Bus, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	return &Bus{conn}, nil
}

func (bus *Bus) SendEvent(interface_name string, event_name string, params ...interface{}) error {
	err := bus.dbus_conn.Emit(path_root, fmt.Sprintf("%s.%s", interface_name, event_name), params...)
	return err
}

func (bus *Bus) RegisterForEvent(interface_name string, event_name string) error {
	filter := make([]string, 2)
	filter[0] = "type='signal'"
	filter[1] = fmt.Sprintf("path='%s'", string(path_root))
	if interface_name != "" {
		filter = append(filter, fmt.Sprintf("interface='%s'", interface_name))
	}
	if event_name != "" {
		filter = append(filter, fmt.Sprintf("member='%s'", event_name))
	}
	call := bus.dbus_conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, strings.Join(filter, ","))
	return call.Err
}

func (bus *Bus) GetEventChannel(size int) chan *dbus.Signal {
	ch := make(chan *dbus.Signal, size)
	bus.dbus_conn.Signal(ch)
	return ch
}
