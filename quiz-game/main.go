package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	csvFilePath = flag.String("file", "./problems.csv", "CSV file path containing quiz questions and answers")
	timeLimit   = flag.Int("limit", 30, "Quiz time limit (in seconds)")
)

type quiz struct {
	questions      int
	answers        int
	correctAnswers int
}

func main() {
	flag.Parse()

	if !startQuiz() {
		fmt.Println("Quiz exited!")
		return
	}
	records, quiz := loadQuiz()
	msg := runQuiz(records, quiz)
	reportQuiz(msg, quiz)
}

func startQuiz() bool {
	fmt.Println("Start quiz (y/n)? ")

	var input string
	fmt.Scan(&input)

	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

func loadQuiz() ([][]string, *quiz) {
	f, err := os.Open(*csvFilePath)
	if err != nil {
		log.Fatal("Failed to open quiz CSV file: ", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("Failed to parse quiz CSV file: ", err)
	}

	return records, &quiz{questions: len(records)}
}

func runQuiz(records [][]string, quiz *quiz) string {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, record := range records {
		question, correctAnswer := strings.TrimSpace(record[0]), strings.TrimSpace(record[1])
		fmt.Printf("Question %d: %v\n", i, question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scan(&answer)
			answerCh <- answer
		}()

		select {
		case answer := <-answerCh:
			fmt.Println()
			quiz.answers++
			if strings.EqualFold(strings.TrimSpace(answer), correctAnswer) {
				quiz.correctAnswers++
			}
		case <-timer.C:
			return "Quiz timed out!"
		}
	}
	return "Quiz completed!"
}

func reportQuiz(msg string, quiz *quiz) {
	fmt.Println(msg)
	fmt.Printf("You answered %d questions and got %d correct out of %d questions", quiz.answers, quiz.correctAnswers, quiz.questions)
}
