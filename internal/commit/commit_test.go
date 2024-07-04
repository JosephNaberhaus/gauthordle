package commit

import (
	"testing"

	"github.com/josephnaberhaus/gauthordle/internal/config"
	"github.com/josephnaberhaus/gauthordle/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_Filter_FilterOutBots(t *testing.T) {
	input := []git.Commit{
		{
			AuthorName:  "Real person",
			AuthorEmail: "joe.smith@example.com",
		},
		{
			AuthorName:  "No reply",
			AuthorEmail: "joe.smith@example.noreply.com",
		},
		{
			AuthorName:  "Robot",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	expected := []git.Commit{
		{
			AuthorName:  "Real person",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	filter, err := BuildFilter()
	require.NoError(t, err)

	assert.Equal(t, expected, filter.filterOutBots(input))
}

func TestFilter_Filter_FiltersExclusions(t *testing.T) {
	input := []git.Commit{
		{
			AuthorName:  "Person 1",
			AuthorEmail: "abc@abc.com",
		},
		{
			AuthorName:  "Person 2",
			AuthorEmail: "def@def.com",
		},
		{
			AuthorName:  "Person 3",
			AuthorEmail: "ghi@ghi.com",
		},
	}

	tests := []struct {
		desc string
		cfg  config.Config
		exp  []git.Commit
	}{{
		desc: "empty config should keep all",
		cfg:  config.Config{},
		exp:  input,
	}, {
		desc: "filter by name",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "2",
			}},
		},
		exp: []git.Commit{input[0], input[2]},
	}, {
		desc: "filter by email",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeEmail: "ghi",
			}},
		},
		exp: []git.Commit{input[0], input[1]},
	}, {
		desc: "specify name and email, prefer name",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName:  "2",
				ExcludeEmail: "ghi",
			}},
		},
		exp: []git.Commit{input[0], input[2]},
	}, {
		desc: "specify multiple filters",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "2",
			}, {
				ExcludeEmail: "ghi",
			}},
		},
		exp: []git.Commit{input[0]},
	}, {
		desc: "filter by name regexp",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeName: "[2-3]",
			}},
		},
		exp: []git.Commit{input[0]},
	}, {
		desc: "filter by email regexp",
		cfg: config.Config{
			AuthorFilters: []config.AuthorFilter{{
				ExcludeEmail: "(abc|def)",
			}},
		},
		exp: []git.Commit{input[2]},
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			filter, err := BuildFilter(WithConfig(tc.cfg))
			require.NoError(t, err)

			assert.Equal(t, tc.exp, filter.filterExclusions(input))
		})
	}
}

func TestFilter_Filter_FilterTeam(t *testing.T) {
	input := []git.Commit{
		{
			AuthorName:  "Person 1",
			AuthorEmail: "abc@abc.com",
		},
		{
			AuthorName:  "Person 2",
			AuthorEmail: "def@def.com",
		},
		{
			AuthorName:  "Person 3",
			AuthorEmail: "ghi@ghi.com",
		},
	}

	testCases := []struct {
		name     string
		cfg      config.Config
		team     string
		expected []git.Commit
	}{
		{
			name:     "invalid team",
			team:     "",
			expected: input,
		},
		{
			name: "unused team",
			cfg: config.Config{
				Teams: map[string]config.Team{
					"rocket": {
						"abc@abc.com",
					},
				},
			},
			expected: input,
		},
		{
			name: "valid team",
			cfg: config.Config{
				Teams: map[string]config.Team{
					"rocket": {
						"abc@abc.com",
					},
				},
			},
			team:     "rocket",
			expected: input[:1],
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			filter, err := BuildFilter(
				WithConfig(testCase.cfg),
				WithTeam(testCase.team),
			)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, filter.filterByTeam(input))
		})
	}
}

func TestFilter_Filter_ConsolidatesAuthorDetails(t *testing.T) {
	input := []git.Commit{
		{
			AuthorName:  "Joe Smith",
			AuthorEmail: "jOe.smith@example.com",
		},
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	expected := []git.Commit{
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
		{
			AuthorName:  "Joseph Smith",
			AuthorEmail: "joe.smith@example.com",
		},
	}

	filter, err := BuildFilter()
	require.NoError(t, err)

	assert.Equal(t, expected, filter.consolidateAuthorDetails(input))
}

func TestFilter_Filter_FiltersCommitSubjects(t *testing.T) {
	input := []git.Commit{
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "DUPLICATED commit message",
		},
		{
			AuthorEmail: "bob.barker@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			SubjectLine: "oneword",
		},
		{
			SubjectLine: "two words",
		},
		{
			SubjectLine: "three words now",
		},
		{
			SubjectLine: "Merge branch 'master' of...",
		},
		{
			SubjectLine: "Merged master",
		},
		{
			SubjectLine: "PREFIX Merge branch 'master' of...",
		},
	}

	expected := []git.Commit{
		{
			AuthorEmail: "joe.smith@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			AuthorEmail: "bob.barker@example.com",
			SubjectLine: "duplicated commit message",
		},
		{
			SubjectLine: "three words now",
		},
	}

	filter, err := BuildFilter()
	require.NoError(t, err)

	assert.Equal(t, expected, filter.filterCommitSubjects(input))
}
