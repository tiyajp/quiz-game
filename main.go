package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	rightAns       = 0
	totalQuestions = 0
	ch             = make(chan []string)
)

func startTimer(start string, myTime *int) *time.Timer {
	if start == "y" {
		return time.NewTimer(time.Duration(*myTime) * time.Second)
	}
	return time.NewTimer(0)
}

func processCSV(rc io.Reader) [][]string {
	r := csv.NewReader(rc)
	rec, err := r.ReadAll()
	totalQuestions = len(rec)
	if err != nil {
		log.Println(err)
		return rec
	}
	return rec
}

func getRecords(file *os.File) {
	rec := processCSV(file)
	for _, records := range rec {
		ch <- records
	}
}

func checkTime(timer *time.Timer) {
	select {
	case <-timer.C:
		fmt.Println("\nYour timer expired!!!")
		fmt.Printf("Score: %d/%d\n", rightAns, totalQuestions)
		close(ch)
		os.Exit(1)
	}
}

func main() {
	problemCount := 0
	quizfile := flag.String("quizFile", "problems.csv", "csv quizFile file in the format of 'question,answer' ")
	myTime := flag.Int("timer", 30, "time to complete the quiz")
	flag.Parse()

	defer close(ch)

	file, err := os.Open(*quizfile)
	if err != nil {
		log.Println(err)
		return
	}

	intro := fmt.Sprintf("Press y to start your timer for %d seconds", *myTime)
	fmt.Println(intro)
	var start string
	fmt.Scan(&start)
	timer := startTimer(start, myTime)

	go getRecords(file)
	for {
		select {
		case rec := <-ch:
			go checkTime(timer)
			var ans string
			problemCount++
			fmt.Printf("Problem #%d: %s\n", problemCount, rec[0])
			fmt.Scan(&ans)
			if ans == strings.TrimSpace(rec[1]) {
				rightAns++
			}

		}
	}

}