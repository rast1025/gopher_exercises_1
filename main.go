package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Quiz struct {
	Question string
	Answer   string
}

func readCSV(filename string) []Quiz {
	csvFile, err := os.Open(filename)
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

func startQuiz(quizSlice []Quiz) {
	correctCount := 0
	fmt.Printf("Quiz started! Total question count: %d.\n", len(quizSlice))
	for i, el := range quizSlice {
		var userAnswer string
		fmt.Printf("Question %d. \n%s = ", i+1, el.Question)
		_, err := fmt.Scanln(&userAnswer)
		if err != nil {
			fmt.Printf("Could not get input from user. %s\n", err)
		}
		if userAnswer == el.Answer {
			correctCount++
		}

	}
	fmt.Printf("Total correct %d out of %d \n", correctCount, len(quizSlice))
}

func main() {
	filename := flag.String("f", "problems.csv", "csv file to parse quiz from")
	flag.Parse()
	quizSlice := readCSV(*filename)
	startQuiz(quizSlice)
}
