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
	"os/user"
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

func timer(mode string, length time.Duration, inputByte chan bool) {
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

func pomodoro(length time.Duration, inputByte chan bool) {
	timer("Pomodoro", length, inputByte)
	clearLine()
}

func do_break(length time.Duration, inputByte chan bool) {
	fmt.Printf("Press any key to start break")
	anyKey(inputByte)
	timer("Break", length, inputByte)
	clearLine()
	fmt.Printf("Break complete. Press any key to start Pomodoro")
	anyKey(inputByte)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// This is probably not valid CSV
// It's just meant to be human-readable
func log_pomodoro(name string, start time.Time, end time.Time, length int) {
	var log_file *os.File
	user, err := user.Current()
	check(err)
	filename := user.HomeDir + "/pomodoros.csv"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log_file, err = os.Create(filename)
		check(err)
		log_file.WriteString("task,timeStart,timeEnd,length\n")
	} else {
		log_file, err = os.OpenFile(filename, os.O_APPEND | os.O_WRONLY, 0600)
		check(err)
	}
	format := "Jan 2 15:04:05"
	log_file.WriteString(fmt.Sprintf("%s,%s,%s,%d\n", name, start.Format(format), end.Format(format), length))
	log_file.Sync()
	log_file.Close()
}

func main() {
	dontDisplayInput()
	var minuteDuration int
	var breakDuration int
	var rounds int
	var task_name string
	flag.IntVar(&minuteDuration, "p", 25, "Duration of pomodoro in minutes")
	flag.IntVar(&breakDuration, "b", 5, "Duration of breaks in minutes")
	flag.IntVar(&rounds, "r", 4, "Amount of pomodoros")
	flag.StringVar(&task_name, "t", "", "Give a name to a task to log it in ~/pomodoros.csv")
	flag.Parse()

	var inputByte = make(chan bool, 100)
	go readStdin(inputByte)

	for i := 1; i <= rounds; i++ {
		startTime := time.Now()
		pomodoro(time.Minute*time.Duration(minuteDuration), inputByte)
		fmt.Printf("Pomodoro %d/%d completed\n", i, rounds)
		if (task_name != "") {
			log_pomodoro(task_name, startTime, time.Now(), minuteDuration)
		}
		if i != rounds {
			do_break(time.Minute*time.Duration(breakDuration), inputByte)
		}
	}
}
