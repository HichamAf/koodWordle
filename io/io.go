package io

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// LoadWords reads a word list from a file.
func LoadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words, scanner.Err()
}

// SaveStats appends a game record to CSV.
func SaveStats(filename, username, secret string, attempts int, status string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error saving stats:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{username, secret, strconv.Itoa(attempts), status}
	writer.Write(record)
}

type Stats struct {
	User          string
	GamesPlayed   int
	GamesWon      int
	TotalAttempts int
}

// LoadStats loads stats for a given user.
func LoadStats(filename, username string) (Stats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Stats{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return Stats{}, err
	}

	stats := Stats{User: username}
	for _, r := range records {
		if r[0] == username {
			stats.GamesPlayed++
			if r[3] == "win" {
				stats.GamesWon++
			}
			att, _ := strconv.Atoi(r[2])
			stats.TotalAttempts += att
		}
	}
	return stats, nil
}

// Print prints the user stats to stdout.
func (s Stats) Print() {
	avg := 0.0
	if s.GamesPlayed > 0 {
		avg = float64(s.TotalAttempts) / float64(s.GamesPlayed)
	}
	fmt.Printf("Stats for %s:\n", s.User)
	fmt.Println("Games played:", s.GamesPlayed)
	fmt.Println("Games won:", s.GamesWon)
	fmt.Printf("Average attempts per game: %.2f\n", avg)
}
