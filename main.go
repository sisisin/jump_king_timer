package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <start_time>")
		return
	}

	ticker := time.NewTicker(1 * time.Second)

	zero, _ := time.Parse("15:04:05", "00:00:00")
	inputStartTime, err := time.Parse("15:04:05", os.Args[1])
	if err != nil {
		err := fmt.Errorf("invalid time format: %s, use hh:mm:ss", os.Args[1])
		panic(err)
	}

	startTime := inputStartTime.Sub(zero)
	now := time.Now()

	fmt.Print("\033[1;1H\033[2J")
	fmt.Print("Timer: ")
	write(startTime)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done")
			return
		case <-ticker.C:
			elapsed := startTime + time.Since(now)

			write(elapsed)
		}
	}
}

const (
	windShiftDuration = 6.527979824107605
)

var n = 0
var val = 5.5

func calcWindShiftN(sec float64) float64 {
	for {
		if sec < val {
			break
		}
		val = 5.5 + windShiftDuration*float64(n)
		n++
	}
	return val
}

func write(elapsed time.Duration) {
	sec := elapsed.Seconds()
	windDirection := ""
	if sec <= 5.5 {
		windDirection = "右"
	} else if int((sec-5.5)/windShiftDuration)%2 == 0 {
		windDirection = "左"
	} else {
		windDirection = "右"
	}

	fmt.Print("\033[?25l")       // カーソルを非表示にする
	defer fmt.Print("\033[?25h") // 関数終了時にカーソルを再表示する

	fmt.Print("\033[1;8H\033[K")
	fmt.Printf("%s\n", formatDuration(elapsed))
	fmt.Printf("%s\n", windDirection)
	fmt.Printf("next: %s\n", formatDuration(time.Duration(calcWindShiftN(sec)*float64(time.Second))))
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
