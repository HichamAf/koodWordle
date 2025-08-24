package game

import (
	"bufio"
	"fmt"
	"strings"
)

const maxAttempts = 6

// ANSI colors
const (
	green  = "\u001B[32m"
	yellow = "\u001B[33m"
	white  = "\u001B[37m"
	reset  = "\u001B[0m"
)

// Play starts the Wordle game.
func Play(scanner *bufio.Scanner, secret string, wordList []string) (bool, int) {
	fmt.Println("Welcome to Wordle! Guess the 5-letter word.")

	attempts := 0
	secret = strings.ToUpper(secret)

	remainingLetters := make(map[rune]bool)
	for r := 'A'; r <= 'Z'; r++ {
		remainingLetters[r] = true
	}

	for attempts < maxAttempts {
		fmt.Print("Enter your guess:  ")
		if !scanner.Scan() {
			break
		}
		guess := scanner.Text()

		if len(guess) != 5 {
			fmt.Println("Your guess must be exactly 5 letters long.")
			continue
		}
		if !isLowercase(guess) {
			fmt.Println("Your guess must only contain lowercase letters.")
			continue
		}
		if !isWordInList(guess, wordList) {
			fmt.Println("Word not in list. Please enter a valid word.")
			continue
		}

		attempts++

		if strings.ToUpper(guess) == secret {
			fmt.Println("Congratulations! You've guessed the word correctly.")
			fmt.Print("Do you want to see your stats? (yes/no): ")
			return true, attempts
		}

		feedback := generateFeedback(secret, guess)
		fmt.Println("Feedback:", feedback)

		// Update remaining letters
		for _, ch := range strings.ToUpper(guess) {
			if !strings.ContainsRune(secret, ch) {
				remainingLetters[ch] = false
			}
		}

		// Print remaining letters
		fmt.Print("Remaining letters: ")
		for r := 'A'; r <= 'Z'; r++ {
			if remainingLetters[r] {
				fmt.Printf("%c ", r)
			}
		}
		fmt.Println()

		// Print remaining attempts
		fmt.Println("Attempts remaining: ", maxAttempts-attempts)

	}

	// Only after the loop ends
	if attempts >= maxAttempts {
		fmt.Printf("Game over. The correct word was: %s\n", strings.ToLower(secret))
		fmt.Print("Do you want to see your stats? (yes/no): ")
	}
	return false, attempts
}

func generateFeedback(secret, guess string) string {
	secret = strings.ToUpper(secret)
	guess = strings.ToUpper(guess)
	result := ""

	for i := 0; i < 5; i++ {
		if guess[i] == secret[i] {
			result += green + string(guess[i]) + reset
		} else if strings.ContainsRune(secret, rune(guess[i])) {
			result += yellow + string(guess[i]) + reset
		} else {
			result += white + string(guess[i]) + reset
		}
	}
	return result
}

func isLowercase(s string) bool {
	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}

func isWordInList(word string, list []string) bool {
	for _, w := range list {
		if w == word {
			return true
		}
	}
	return false
}
