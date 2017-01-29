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
	"os"
	"os/signal"
	"syscall"
	"time"
)

const weather_poll_sec int = 60

func update_weather() {
	fmt.Println("Updating weather...")
}

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second * time.Duration(weather_poll_sec))
	go func() {
		for {
			update_weather()
			<-ticker.C
		}
	}()

	select {
	case sig := <-ch:
		// TODO log to syslog
		fmt.Println("Received signal", sig)
	}
}
