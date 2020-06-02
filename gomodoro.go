/*
 * This file is part of the gomodoro application
 * Copyright (c) 2020 Rasmus Rendal.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func clearLine() {
	fmt.Printf("\r")
	fmt.Printf(strings.Repeat(" ", 50))
	fmt.Printf("\r")
	//fmt.Printf(strings.Repeat("\b", 40))
}

func anyKey(inputByte chan bool) {
	_ = <-inputByte
}

func readStdin(x chan bool) {
	for true {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		if b[0] == 32 {
			x <- true
		} else {
			x <- false
		}
	}
}

func pomodoro(mode string, length time.Duration, inputByte chan bool) {
	endTime := time.Now().Add(length)

	for time.Now().Before(endTime) {
		timeLeft := endTime.Sub(time.Now())
		m := timeLeft / time.Minute
		timeLeft -= m * time.Minute
		s := timeLeft / time.Second
		clearLine()
		fmt.Printf(mode+" time remaining: %02d:%02d", m, s)
		if len(inputByte) > 0 {
			var x = <-inputByte
			if x {
				var sleepStartTime = time.Now()
				fmt.Printf(" - Paused")
				x = <-inputByte
				for !x {
					x = <-inputByte
				}
				endTime = endTime.Add(time.Now().Sub(sleepStartTime))
			}
		}
		time.Sleep(time.Second)
	}
}

func dontDisplayInput() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func main() {
	dontDisplayInput()
	var minuteDuration int
	var breakDuration int
	var rounds int
	flag.IntVar(&minuteDuration, "p", 25, "Duration of pomodoro in minutes")
	flag.IntVar(&breakDuration, "b", 5, "Duration of breaks in minutes")
	flag.IntVar(&rounds, "r", 4, "Amount of pomodoros")
	flag.Parse()

	var inputByte = make(chan bool, 100)
	go readStdin(inputByte)

	for i := 1; i <= rounds; i++ {
		pomodoro("Pomodoro", time.Minute*time.Duration(minuteDuration), inputByte)
		clearLine()
		fmt.Printf("Pomodoro %d/%d completed\n", i, rounds)
		if i != rounds {
			fmt.Printf("Press any key to start break")
			anyKey(inputByte)
			pomodoro("Break", time.Minute*time.Duration(breakDuration), inputByte)
			clearLine()
			fmt.Printf("Break complete. Press any key to start Pomodoro")
			anyKey(inputByte)
		}
	}

}
