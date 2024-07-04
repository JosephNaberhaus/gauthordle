package game

import (
	"github.com/JosephNaberhaus/prompt"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	"github.com/josephnaberhaus/gauthordle/internal/output"
	"strconv"
	"time"
)

const header = "                   _   _                   _ _      \n  __ _  __ _ _   _| |_| |__   ___  _ __ __| | | ___ \n / _` |/ _` | | | | __| '_ \\ / _ \\| '__/ _` | |/ _ \\\n| (_| | (_| | |_| | |_| | | | (_) | | | (_| | |  __/\n \\__, |\\__,_|\\__,_|\\__|_| |_|\\___/|_|  \\__,_|_|\\___|\n |___/      The daily git author guessing game      "
const youWin = "                              _       \n _   _  ___  _   _  __      _(_)_ __  \n| | | |/ _ \\| | | | \\ \\ /\\ / / | '_ \\ \n| |_| | (_) | |_| |  \\ V  V /| | | | |\n \\__, |\\___/ \\__,_|   \\_/\\_/ |_|_| |_|\n |___/                                "
const youLose = "                     _                \n _   _  ___  _   _  | | ___  ___  ___ \n| | | |/ _ \\| | | | | |/ _ \\/ __|/ _ \\\n| |_| | (_) | |_| | | | (_) \\__ \\  __/\n \\__, |\\___/ \\__,_| |_|\\___/|___/\\___|\n |___/                                "
const nope = "                        \n _ __   ___  _ __   ___ \n| '_ \\ / _ \\| '_ \\ / _ \\\n| | | | (_) | |_) |  __/\n|_| |_|\\___/| .__/ \\___|\n            |_|         "

const numPuzzleCommits = 4

type puzzleHints struct {
	totalCommits    int
	mostTouchedFile string
}

type Puzzle struct {
	authorEmail   string
	authorName    string
	authorCommits []git.Commit
	puzzleCommits [numPuzzleCommits]git.Commit

	hints puzzleHints

	// All commits by all users.
	allCommits     []git.Commit
	allAuthorNames map[string]string
}

func (p Puzzle) Run() error {
	var promptOptions []prompt.SelectionOption
	for authorEmail, authorName := range p.allAuthorNames {
		promptOptions = append(promptOptions, prompt.SelectionOption{
			ID:          authorEmail,
			Name:        authorName,
			Description: authorEmail,
		})
	}

	for stage := 0; stage < numPuzzleCommits; stage++ {
		output.ClearScreen()
		output.PrintColorLn(header, output.Yellow)
		output.Ln()

		output.PrintColor("Guess the author of the following commit", output.White)
		if stage > 0 {
			output.PrintColor("s", output.White)
		}
		output.PrintColorLn(":", output.Yellow)
		output.Ln()

		for i := 0; i <= stage; i++ {
			output.PrintColor("Commit #", output.Green)
			output.PrintColor(strconv.Itoa(i+1), output.Green)
			output.PrintColor(": ", output.Green)
			output.PrintColorLn(p.puzzleCommits[i].SubjectLine, output.White)
		}

		// Hints
		output.Ln()
		if stage >= 1 {
			output.Ln()
			output.PrintColorLn("Hints", output.Green)
			output.PrintColor("Number of commits made by author in the last year: ", output.Green)
			output.PrintColorLn(strconv.Itoa(p.hints.totalCommits), output.White)
		}
		if stage >= 3 {
			output.PrintColor("Author's most touched file: ", output.Green)
			output.PrintColorLn(p.hints.mostTouchedFile, output.White)
		}

		output.Ln()
		answerPrompt := &prompt.Select{
			Question:      "Who is the author?",
			Options:       promptOptions,
			NumLinesShown: 3,
		}
		err := answerPrompt.Show()
		if err != nil {
			return err
		}

		if answerPrompt.Response().ID == p.authorEmail {
			flashMessage(youWin, output.Green)
			break
		} else {
			if stage == numPuzzleCommits-1 {
				flashMessage(youLose, output.Red)
			} else {
				flashMessage(nope, output.Red)
			}
		}
	}

	output.Ln()
	output.PrintColor("The answer was: ", output.White)
	output.PrintColor(p.authorName, output.White)
	output.PrintColor(" (", output.White)
	output.PrintColor(p.authorEmail, output.White)
	output.PrintColorLn(")", output.White)
	output.Ln()

	return nil
}

func flashMessage(message string, color output.Color) {
	output.ClearScreen()
	output.PrintColorLn(header, output.Yellow)
	output.Ln()
	output.PrintColorLn(message, color)
	time.Sleep(1000 * time.Millisecond)
}
