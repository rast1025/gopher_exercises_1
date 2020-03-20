package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	Question string
	Answer   string
}

func readCSV(filename string) []Quiz {
	csvFile, err := os.Open(filename)
	defer func() {
		err := csvFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	quizSlice := []Quiz{}
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(csvFile)
	for {
		quiz := Quiz{}
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(data) != 2 {
			log.Fatal("Csv file must consist of 2 strings")
		}
		quiz.Question = data[0]
		quiz.Answer = data[1]
		quizSlice = append(quizSlice, quiz)
	}
	if len(quizSlice) == 0 {
		log.Fatal("No questions found in csv file")
	}
	return quizSlice
}

func startQuiz(quizSlice []Quiz, duration uint) {
	fmt.Printf("Time given for quiz is %d second(s). Press enter to start the quiz\n", duration)
	_, err := fmt.Scanln()
	if err != nil {
		log.Fatal(err)
	}
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	correctCount := 0
	fmt.Printf("Quiz started! Total question count: %d.\n", len(quizSlice))
	for i, el := range quizSlice {
		fmt.Printf("Question %d. \n%s = ", i+1, el.Question)
		answersCh := make(chan string)
		go func() {
			var userAnswer string
			_, err := fmt.Scanln(&userAnswer)
			if err != nil {
				log.Fatal("Could not get input from user. %s\n", err)
			}
			answersCh <- userAnswer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTotal correct %d out of %d \n", correctCount, len(quizSlice))
			return
		case userAnswer := <-answersCh:
			if strings.ToLower(strings.TrimSpace(userAnswer)) == strings.ToLower(el.Answer) {
				correctCount++
			}
		}
	}
	fmt.Printf("Total correct %d out of %d \n", correctCount, len(quizSlice))

}

func main() {
	filename := flag.String("f", "problems.csv", "csv file to parse quiz from")
	duration := flag.Uint("d", 30, "duration of the quiz")
	shuffle := flag.Bool("s", false, "whether to shuffle quiz questions")
	flag.Parse()
	quizSlice := readCSV(*filename)
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizSlice), func(i, j int) { quizSlice[i], quizSlice[j] = quizSlice[j], quizSlice[i] })
	}
	startQuiz(quizSlice, *duration)
}
