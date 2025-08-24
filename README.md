# koodWordle

A command-line Wordle-like game written in Go. Guess a hidden 5-letter word in up to six attempts. After each guess, the game provides colorized feedback (Green/Yellow/White), shows the number of attempts remaining, and prints the list of remaining (not-yet-confirmed-incorrect) letters A–Z. Game statistics persist across sessions in a CSV file.

---

## Features

- **6 attempts** to guess the secret **5-letter** word
- **Color feedback** using ANSI escape codes:
  - **Green** (`\u001B[32m`) — correct letter in the **correct position**
  - **Yellow** (`\u001B[33m`) — correct letter in the **wrong position**
  - **White** (`\u001B[37m`) — incorrect letter
- **Remaining letters** view (A–Z, uppercase, whitespace separated)
- **Attempts remaining** after every guess
- **Username-based stats** persisted to `stats.csv`:
  - Username, Secret word, Attempts, `win|loss`
- **Graceful handling** of missing/invalid arguments, missing word list, and `EOF` during input
- **Uses** `bufio.Scanner` for line-based stdin reading
- **Testable** via command-line index: `go run . 10`

---

## Game Rules (User Experience)

- You’ll be asked for a **username**.
- You’ll see: `Welcome to Wordle! Guess the 5-letter word.`
- Enter guesses (letters only). After each guess you’ll see:
  - `Feedback: <GUESS IN UPPERCASE WITH COLORED LETTERS>`
  - `Remaining letters: <A B C ... Z>` (excluding letters known to be incorrect)
  - `Attempts remaining: <N>`
- If you **guess the word** within 6 attempts, you win.
- If you **don’t**, the game reveals the **secret word**.
- After the game, you can opt to **view your stats**.

---

## Example Session (Stdout Format)

> The program’s stdout matches these formats exactly to align with tests.

```
Enter your username:
alice

Welcome to Wordle! Guess the 5-letter word.
Enter your guess:
crane
Feedback: CRANE
Remaining letters: A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
Attempts remaining: 5
Enter your guess:
slate
Feedback: SLATE
Remaining letters: A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
Attempts remaining: 4
Enter your guess:
...

Do you want to see your stats? (yes/no):
yes
Stats for alice:
Games played: 3
Games won: 2
Average attempts per game: 4.33
Press Enter to exit...
```

> Note: Actual feedback prints each letter colorized using ANSI codes while the **letters themselves are uppercase**.

---

## Color & Feedback Scheme

- Letters in the feedback are **capitalized**.
- Each letter of the guess is evaluated:
  1. **Green** if it matches the secret letter at the same position.
  2. **Yellow** if it exists elsewhere in the secret word and hasn’t been fully matched already (duplicate-aware).
  3. **White** otherwise.

**ANSI Escape Codes Used**
- Green: `\u001B[32m`
- Yellow: `\u001B[33m`
- White: `\u001B[37m`
- Reset: `\u001B[0m`

---

## Remaining Letters

After each guess, the game prints **Remaining letters** — a **sorted list A–Z** (uppercase, whitespace separated) that includes:
- **Letters not yet guessed incorrectly**, plus letters confirmed present (Green/Yellow).
- **Excludes** letters known to be **incorrect** (accumulated across attempts).

---

## Stats Storage (`stats.csv`)

The file `stats.csv` is created at the project root and **persists across game sessions**. Each game appends one row with **comma-separated values**:

| Column          | Description                          | Example   |
|-----------------|--------------------------------------|-----------|
| Username        | Name entered at login                | `alice`   |
| Secret word     | The word for that game               | `SLATE`   |
| Attempts        | Number of guesses made (1–6)         | `4`       |
| Result          | `win` or `loss`                      | `win`     |

> Both `stats.csv` and `wordle-words.txt` are **ignored by Git** (see `.gitignore`).

---

## Project Structure

```
koodWordle
├── main.go           // Entry point, processes arguments, starts game
├── game/             // Game logic and mechanics
│   └── game.go       // Core game functionality and feedback generation
├── io/               // Input/output operations
│   └── io.go         // File handling for words and statistics
└── model/            // Data structures
    └── user.go       // User entity and statistics tracking
```

- **Module name** matches repository name: `koodWordle`.
- `main.go` is at the **project root**.

---

## Running the Game

### Prerequisites
- Go 1.20+

### Clone and Run
```bash
git clone <your-repo-url> koodWordle
cd koodWordle

# Make sure the module name matches the repo name
go mod init koodWordle || true
go mod tidy
```

Place your **word list** file at project root as `wordle-words.txt`. Do **not** commit it.

Run the game with a **word index** (for testability):
```bash
go run . 10
```

- The program reads `wordle-words.txt`, selects the word at index `10` (0-based or 1-based configurable in code; this implementation uses **0-based** for predictability), and starts the session.

---

## Implementation Notes

### Input & EOF Handling
- Input is read using `bufio.Scanner`.
- `scanner.Scan()` returning `false` (including `EOF` from `Ctrl+D`) **breaks loops gracefully** without crashing.

### Word List
- Read from `wordle-words.txt` (UTF-8, one word per line, 5 letters).
- Missing file is handled gracefully with a clear message and exit.

### Command-Line Argument
- Index is provided as the **first argument** (e.g., `go run . 10`).
- Missing/invalid index is handled gracefully:
  - Prints a helpful message (`Usage: go run . <word-index>`) and exits with a non-zero status **without crashing**.

### Feedback Generation
- Implements a **two-pass** algorithm to handle duplicates correctly:
  1. First pass marks **Green** matches and counts remaining letters in the secret.
  2. Second pass assigns **Yellow** for remaining instances; otherwise **White**.

### Remaining Letters Tracking
- Maintains a set of **letters marked incorrect** across attempts.
- Renders remaining letters as `A B C ... Z` minus incorrect letters.

### Stats
- Appends to `stats.csv` after each game.
- On **stats view**, aggregates:
  - `Games played`
  - `Games won`
  - `Average attempts per game` (float with two decimals)

### ANSI Output
- Colors are applied per-letter using escape codes; a final **reset** is emitted after each colored token.

---

## .gitignore

```
wordle-words.txt
stats.csv
```

> Ensure both files are **excluded** from the repository.

---

## Testing Hints

- Use deterministic indices (`go run . 0`, `go run . 1`, …) to pick known words in `wordle-words.txt`.
- Validate stdout **exact wording**:
  - `Enter your username:`
  - `Welcome to Wordle! Guess the 5-letter word.`
  - `Enter your guess:`
  - `Feedback: <UPPERCASE GUESS>`
  - `Remaining letters: <LETTERS>`
  - `Attempts remaining: <N>`
  - `Do you want to see your stats? (yes/no):`
  - `Stats for <username>:`
  - `Games played: <number>`
  - `Games won: <number>`
  - `Average attempts per game: <float number>`
  - `Press Enter to exit...`

---

## Example Edge Cases

- **Missing argument**: prints usage and exits.
- **Invalid index**: prints an informative error and exits.
- **Missing `wordle-words.txt`**: prints a friendly error and exits.
- **Non-5-letter guess**: prompts user again without consuming an attempt.
- **Non-alpha input**: ignored with a hint; doesn’t consume an attempt.
- **EOF during username or guesses**: exits gracefully.

---

## Build

```bash
go build -o koodWordle
./koodWordle 10
```

---

## Future Improvements

- Daily word mode
- Hard mode (must reuse discovered letters)
- Colorblind-friendly symbols
- Import/export stats
- In-game help (`?` to show rules)