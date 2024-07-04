package game

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/josephnaberhaus/gauthordle/internal/git"
)

func commitsByAuthorEmail(commits []git.Commit) map[string][]git.Commit {
	result := map[string][]git.Commit{}
	for _, commit := range commits {
		result[commit.AuthorEmail] = append(result[commit.AuthorEmail], commit)
	}

	return result
}

func numCommitsByAuthorEmail(commits []git.Commit) map[string]int {
	result := map[string]int{}
	for authorEmail, commits := range commitsByAuthorEmail(commits) {
		result[authorEmail] = len(commits)
	}

	return result
}

func nameByEmail(commits []git.Commit) map[string]string {
	result := map[string]string{}
	for _, commit := range commits {
		result[commit.AuthorEmail] = commit.AuthorName
	}

	return result
}

func allAuthorEmails(commits []git.Commit) []string {
	authorSet := map[string]struct{}{}
	for _, commit := range commits {
		authorSet[commit.AuthorEmail] = struct{}{}
	}

	result := make([]string, 0, len(authorSet))
	for author := range authorSet {
		result = append(result, author)
	}

	return result
}

func pickAuthor(commits []git.Commit, authorBias float64, random *rand.Rand) (string, error) {
	numByAuthor := numCommitsByAuthorEmail(commits)
	allAuthors := allAuthorEmails(commits)

	// Filter out authors that have made fewer commits than the number we need for the puzzle.
	allAuthors = slices.DeleteFunc(allAuthors, func(s string) bool {
		if numByAuthor[s] < numPuzzleCommits {
			return true
		}

		return false
	})
	if len(allAuthors) == 0 {
		return "", fmt.Errorf("there are no authors with %d or more valid commits", numPuzzleCommits)
	}

	// Sort the authors by how many commits they've made.
	slices.SortFunc(allAuthors, func(a, b string) int {
		// If the number of commits are the same then just sort by the e-mail.
		if numByAuthor[a] == numByAuthor[b] {
			return strings.Compare(a, b)
		}

		return numByAuthor[a] - numByAuthor[b]
	})

	// Use a root curve so that we favor the higher contributing users.
	// A higher bias increases the likelihood that a high-contributing user will be picked.
	randMax := math.Pow(float64(len(allAuthors)), authorBias)
	randNumber := random.Intn(int(math.Floor(randMax)))
	index := int(math.Pow(float64(randNumber), 1/authorBias))

	// Make sure that floating point error didn't put us in an invalid index
	index = max(index, 0)
	index = min(index, len(allAuthors)-1)

	return allAuthors[index], nil
}

func mostTouchedFileForAuthor(authorEmail string) (string, error) {
	filesChanged, err := git.GetFilesChangedForAuthor(authorEmail)
	if err != nil {
		return "", fmt.Errorf("error while getting the author's most touched file: %w", err)
	}

	fileCount := map[string]int{}
	for _, file := range filesChanged {
		fileCount[file]++
	}

	maxCountFile := ""
	for file, count := range fileCount {
		if count > fileCount[maxCountFile] {
			maxCountFile = file
		}
	}

	return maxCountFile, nil
}
