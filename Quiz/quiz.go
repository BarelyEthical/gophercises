package main

import ("fmt"
	"encoding/csv"
	"log"
	"os"
	"io"
	"flag"
	"time"
)

type Entry struct {
	leftOp, ans string
}

type Score struct {
	correct, incorrect int
}

func readQuiz(f io.Reader, entries *[]Entry) {
	questions, err := csv.NewReader(f).ReadAll()
	if err == io.EOF {
		fmt.Println("Congratulations")
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, items := range questions {
		i := Entry{leftOp : items[0], ans : items[1]}
		*entries = append(*entries, i)
	}
}

func playQuiz(timer *time.Timer, entries []Entry)(s Score) {
	// Read the quiz to the user
	fmt.Println("interactive mode (%v questions in quiz)\n", len(entries))
	for i, items := range entries {
		var op = ""
		// chan to write the user input and then read it in separate goroutine
		c := make(chan string)
		go func() {
			fmt.Printf("%v) %v?: ", i+1, items.leftOp)
			fmt.Scanln(&op)
			c<-op
		}()
		select {
		case <- timer.C:
			fmt.Println("Time out")
			return s
		// wait for read from c chan
		case answer := <-c:
			if answer == items.ans  {
				fmt.Println("correct")
				s.correct++
			} else {
				fmt.Println("incorrect")
				s.incorrect++
			}
		}
	}
	return s
}

func main() {
	// Package OS provides a platform-independent interface to operating system
	// functionality
	f, err := os.Open("problems.csv")
	if err != nil {
		return
	}
	var entries []Entry
	// Package csv reads and writes CSV files.
	readQuiz(f, &entries)

	//Read the timelimit given as the input
	timeLimit := flag.Int("limit",  30,  "timelimit to finish the quiz")
	flag.Parse()
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	var s Score
	s = playQuiz(timer, entries)
	fmt.Printf("Your score is : %v", s.correct)
}
