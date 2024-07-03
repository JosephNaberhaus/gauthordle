package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/josephnaberhaus/gauthordle/internal/commit"
	"github.com/josephnaberhaus/gauthordle/internal/config"
	"github.com/josephnaberhaus/gauthordle/internal/game"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	"github.com/josephnaberhaus/gauthordle/internal/output"
	"math/rand"
	"os"
	"strings"
)

const helpBody = "A daily game where you try to guess the author of some Git commits.\n\nTo play, simply \"git checkout\" the main development branch of your repository\nand run this program with no arguments.\n\nNew games start at midnight Central Time."

var (
	dumpCommits = flag.String("debugDumpCommits", "", "File to dump JSON containing all commits considered when generating the game.")
	help        = flag.Bool("help", false, "Print the help message.")
	random      = flag.Bool("random", false, "If true, play a random game instead of the daily game.")
)

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		exit(fmt.Errorf("unsupported arguments %q\n", strings.Join(flag.Args(), ",")))
	}

	if *help {
		showUsage()
	}

	if !git.IsGitInstalled() {
		exit(errors.New("git must be installed"))
	}

	if !git.IsInGitRepo() {
		exit(errors.New("must be in a git repository"))
	}

	fmt.Println("Building today's game...")

	startTime, endTime := game.PuzzleTimeRange()
	cfg, err := config.Load()
	exitIfError(err)

	filter, err := commit.BuildFilter(
		commit.WithConfig(cfg),
		commit.WithStartTime(startTime),
		commit.WithEndTime(endTime),
	)
	exitIfError(err)
	commits, err := filter.GetCommits()
	exitIfError(err)

	if *dumpCommits != "" {
		serializedCommits, err := json.MarshalIndent(commits, "", "  ")
		exitIfError(err)

		err = os.WriteFile(*dumpCommits, serializedCommits, os.ModePerm)
		exitIfError(err)
	}

	gameOptions := []game.Option{
		game.WithCommits(commits),
	}
	if !*random {
		// For non-random games, use the startTime as the random source so that it's stable throughout the day.
		gameOptions = append(gameOptions, game.WithRandomSource(rand.NewSource(startTime.Unix())))
	}

	puzzle, err := game.BuildPuzzle(gameOptions...)
	exitIfError(err)

	err = puzzle.Run()
	exitIfError(err)
}

func showUsage() {
	fmt.Println(helpBody)
	flag.Usage()

	os.Exit(0)
}

func exitIfError(err error) {
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	output.FprintColor(os.Stderr, fmt.Sprintf("ERROR: %s", err.Error()), output.Red)
	os.Exit(1)
}
