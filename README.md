# Gomodoro
Gomodoro is a pomodoro timer for the terminal, written in go.
It has support for setting timers with variable length, break lengths, and amount of pomodoros. It can also be paused by pressing space.
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

## Time-tracking
There is support for time-tracking in the application, by using the `-t <task name>` argument. 
It will create or append to a file called `~/pomodoros.csv`, adding pomodoros as you do them

Additionally, there is support for setting the length of your pomodoros to zero.
In this case, the timer will count upwards.
You can still start breaks by pressing spacebar.
But additionally, you have to press ctrl+c to stop the timer, and add the task to the tracker.
