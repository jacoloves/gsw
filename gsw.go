package main

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/mattn/go-tty"
)

var (
	marks                       = []string{"|", "/", "-", "\\"}
	out               io.Writer = os.Stdout
	hour_record       []int64
	minute_record     []int64
	seconds_record    []int64
	mirisecond_record []int64
)

func clear() {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")
	fmt.Print()
}

func time_drow(elapsed time.Duration) {
	var h, m, s, ms int64

	h = mod(elapsed.Hours(), 24)
	m = mod(elapsed.Minutes(), 60)
	s = mod(elapsed.Seconds(), 60)
	ms = mod(float64(elapsed.Nanoseconds())/1000, 100)

	fmt.Fprint(os.Stdout, fmt.Sprintf("%v:%v:%v:%v\r", h, m, s, ms))
}

func mod(val float64, mod int64) int64 {
	raw := big.NewInt(int64(val))
	return raw.Mod(raw, big.NewInt(mod)).Int64()
}

func message_press_ent() {
	fmt.Println()
	fmt.Printf("\r%s", "Press Enter Stop")
	fmt.Print("\x1b[H")
}

func time_cout(start time.Time, cancel chan struct{}) {
	var elapsed time.Duration
	for {
		select {
		case <-cancel:
			return
		case <-time.After(time.Millisecond):
			elapsed = time.Since(start)
			time_drow(elapsed)
		}
	}
}

func tty_cencer(start time.Time) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	cancel := make(chan struct{})

	for {
		go time_cout(start, cancel)
		message_press_ent()
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if r == 13 {
			close(cancel)
			return
		}
	}
}

func stop_process(stop_time time.Time) time.Time {
	fmt.Println()
	message := fmt.Sprintf("\r\x1b[K%s", "🖐stop")
	io.WriteString(out, message)
	drawTimeRecord()
	inputKey()
	return time.Now()
}

func inputKey() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if r == 13 {
			return
		}
	}
}

func timeRecord(stop_time time.Time) {
	var elapsed time.Duration
	var h, m, s, ms int64
	elapsed = time.Since(stop_time)

	h = mod(elapsed.Hours(), 24)
	m = mod(elapsed.Minutes(), 60)
	s = mod(elapsed.Seconds(), 60)
	ms = mod(float64(elapsed.Nanoseconds())/1000, 100)

	hour_record = append(hour_record, h)
	minute_record = append(minute_record, m)
	seconds_record = append(seconds_record, s)
	mirisecond_record = append(mirisecond_record, ms)
}

func drawTimeRecord() {
	for i := 0; i < len(hour_record); i++ {
		fmt.Println()
		fmt.Fprint(os.Stdout, fmt.Sprintf("%v  %v:%v:%v:%v", i+1, hour_record[i], minute_record[i], seconds_record[i], mirisecond_record[i]))
	}
}

func main() {
	start := time.Now()
	//clear()
	for {
		clear()
		tty_cencer(start)
		timeRecord(start)
		start = stop_process(start)
	}
}
