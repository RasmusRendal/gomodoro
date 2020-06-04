# Gomodoro
Gomodoro is a pomodoro timer for the terminal, written in go.
It has support for setting timers with variable length, break lengths, and amount of pomodoros. It can also be paused by pressing space.
If you provide the argument `-t <task name>`, a file will be created or appended to in `~/pomodoros.csv`, which will contain the pomodoros performed
This is the first thing I've ever written in Go, so it's probably not great.

The options are:
```
  -b int
    	Duration of breaks in minutes (default 5)
  -p int
    	Duration of pomodoro in minutes (default 25)
  -r int
    	Amount of pomodoros (default 4)
  -t string
    	Task to log (Default no logging)
```
