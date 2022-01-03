package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvPath := flag.String("csv", "problems.csv", "path to the problems.csv file")
	timeLimit := flag.Int("time-limit", 5, "timer in seconds")
	flag.Parse()

	// Read each problem from `problems.csv`
	f, err := os.Open(*csvPath)
	if err != nil {
		exit(fmt.Sprintf("failed to open file %s, with error: %v\n", *csvPath, err))
	}

	defer f.Close()

	// Read all into memory since the problems.csv file will not be very big
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		exit(err.Error())
	}

	problems := parseLines(lines)

	var startInput string
	for startInput != "s" {
		fmt.Println(startInput)
		fmt.Printf("Press s to start quiz: ")
		fmt.Scanf("%s", &startInput)
	}

	t := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	var correct int
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.Question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-t.C:
			fmt.Printf("\nGot %d/%d questions correct!\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.Answer {
				correct++
			}
		}
	}
	fmt.Printf("\nGot %d/%d questions correct!\n", correct, len(problems))
}

func parseLines(input [][]string) []Problem {
	out := make([]Problem, len(input))
	for i, problemSet := range input {
		out[i] = Problem{Question: problemSet[0], Answer: strings.TrimSpace(problemSet[1])}
	}
	return out
}

type Problem struct {
	Question string
	Answer   string
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}
