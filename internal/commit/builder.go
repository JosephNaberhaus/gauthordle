package commit

import (
	"cmp"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/josephnaberhaus/gauthordle/internal/config"
)

type FilterOption func(filter *Filter) error

func BuildFilter(options ...FilterOption) (*Filter, error) {
	filter := new(Filter)
	for _, opt := range options {
		err := opt(filter)
		if err != nil {
			return nil, fmt.Errorf("error build commit filter: %w", err)
		}
	}

	return filter, nil
}

func WithConfig(cfg config.Config) FilterOption {
	return func(filter *Filter) error {
		for _, authorFilter := range cfg.AuthorFilters {
			r, err := regexp.Compile(
				cmp.Or(authorFilter.ExcludeName, authorFilter.ExcludeEmail),
			)
			if err != nil {
				return err
			}

			if authorFilter.ExcludeName != "" {
				filter.nameFilters = append(filter.nameFilters, r)
			} else {
				filter.emailFilters = append(filter.emailFilters, r)
			}
		}

		filter.teams = make(map[string]map[string]struct{}, len(cfg.Teams))
		for name, team := range cfg.Teams {
			for _, email := range team {
				if _, ok := filter.teams[name]; !ok {
					filter.teams[name] = map[string]struct{}{}
				}

				email = strings.ToLower(email)
				filter.teams[name][email] = struct{}{}
			}
		}

		return nil
	}
}

func WithStartTime(startTime time.Time) FilterOption {
	return func(filter *Filter) error {
		filter.startTime = startTime

		return nil
	}
}

func WithEndTime(endTime time.Time) FilterOption {
	return func(filter *Filter) error {
		filter.endTime = endTime

		return nil
	}
}

func WithTeam(team string) FilterOption {
	return func(filter *Filter) error {
		filter.team = team

		return nil
	}
}
