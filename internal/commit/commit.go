package commit

import (
	"fmt"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	"regexp"
	"slices"
	"strings"
	"time"
)

type Filter struct {
	startTime, endTime time.Time

	nameFilters  []*regexp.Regexp
	emailFilters []*regexp.Regexp
}

func (f *Filter) GetCommits() ([]git.Commit, error) {
	if f.startTime.IsZero() {
		return nil, fmt.Errorf("no start time specified")
	}
	if f.endTime.IsZero() {
		return nil, fmt.Errorf("no end time specified")
	}

	commits, err := git.GetCommits(f.startTime, f.endTime)
	if err != nil {
		return nil, err
	}

	type filterFunc func([]git.Commit) []git.Commit
	filters := []filterFunc{
		f.filterForConfig,
		f.filterOutBots,
		f.filterCommitSubjects,
		f.consolidateAuthorDetails,
	}
	for _, filter := range filters {
		commits = filter(commits)
	}

	return commits, nil
}

// filterBots attempts to filter bot-made commits out.
// It's impossible to cover all cases here, but we'll make a best effort.
func (f *Filter) filterOutBots(commits []git.Commit) []git.Commit {
	var result []git.Commit
	for _, commit := range commits {
		// E-mails with "noreply" in them are usually associated with bots.
		if strings.Contains(commit.AuthorEmail, "noreply") {
			continue
		}

		// If the author's name contains the word "robot" that's a pretty good indication that it's a robot.
		if strings.Contains(strings.ToLower(commit.AuthorName), "robot") {
			continue
		}

		result = append(result, commit)
	}

	return result
}

// filterForConfig filters out commits according to the user's config as defined in the config
// package.
func (f *Filter) filterForConfig(commits []git.Commit) []git.Commit {
	shouldDelete := func(commit git.Commit) bool {
		for _, f := range f.nameFilters {
			if f.MatchString(commit.AuthorName) {
				return true
			}
		}
		for _, f := range f.emailFilters {
			if f.MatchString(commit.AuthorEmail) {
				return true
			}
		}
		return false
	}

	commits = slices.Clone(commits)
	return slices.DeleteFunc(commits, shouldDelete)
}

// filterCommitSubjects ensures that we only use interesting commit subjects
func (f *Filter) filterCommitSubjects(commits []git.Commit) []git.Commit {
	var result []git.Commit
	authorCommitSubjects := map[string]map[string]struct{}{}
	for _, commit := range commits {
		// Commit messages with only 1-2 words don't get you very much information.
		if len(strings.Split(commit.SubjectLine, " ")) < 3 {
			continue
		}

		// Merge commits aren't helpful
		if strings.Contains(commit.SubjectLine, "Merge branch") {
			continue
		}

		// Some authors use the same commit message over and over again.
		// Remove the duplicates
		if authorCommitSubjects[commit.AuthorEmail] == nil {
			authorCommitSubjects[commit.AuthorEmail] = map[string]struct{}{}
		}

		// No reason for this to be case sensitive
		subjectLine := strings.ToLower(commit.SubjectLine)

		authorCommits := authorCommitSubjects[commit.AuthorEmail]
		if _, ok := authorCommits[subjectLine]; ok {
			continue
		}
		authorCommits[subjectLine] = struct{}{}

		result = append(result, commit)
	}

	return result
}

// consolidateAuthorDetails ensures that the e-mails and names for all authors are uniform.
func (f *Filter) consolidateAuthorDetails(commits []git.Commit) []git.Commit {
	result := make([]git.Commit, len(commits))
	emailToAuthorName := map[string]string{}

	// Go backwards so that we use the most recently used name
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]

		// E-mails are not case-sensitive, so lower-case it.
		authorEmail := strings.ToLower(commit.AuthorEmail)
		commit.AuthorEmail = authorEmail

		if _, ok := emailToAuthorName[authorEmail]; !ok {
			// This is the first time we have encountered this e-mail.
			emailToAuthorName[authorEmail] = commit.AuthorName
		} else {
			// We've encountered this e-mail before. Make sure we use that same name.
			commit.AuthorName = emailToAuthorName[authorEmail]
		}

		result[i] = commit
	}

	return result
}
