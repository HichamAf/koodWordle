package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"koodWordle/game"
	"koodWordle/io"
	"koodWordle/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a number as command line argument")
		bufio.NewScanner(os.Stdin).Scan()
		return
	}

	index, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid command-line argument. Please launch with a valid number.")
		return
	}

	words, err := io.LoadWords("wordle-words.txt")
	if err != nil {
		fmt.Println("Error loading word list:", err)
		bufio.NewScanner(os.Stdin).Scan()
		return
	}

	if index >= len(words) || index < 0 {
		fmt.Print("Enter your username: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() // consume the 'user' input
		fmt.Println("Invalid word number.")
		fmt.Println("Press Enter to exit...")
		scanner.Scan() // final wait
		return
	}

	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	username := scanner.Text()
	user := model.NewUser(username)

	result, attempts := game.Play(scanner, words[index], words)

	status := "loss"
	if result {
		status = "win"
	}
	io.SaveStats("stats.csv", user.Name, words[index], attempts, status)

	if scanner.Scan() && scanner.Text() == "yes" {
		stats, err := io.LoadStats("stats.csv", user.Name)
		if err == nil {
			stats.Print()
			fmt.Println("Press Enter to exit...")
			bufio.NewScanner(os.Stdin).Scan()
		}
	}
}
