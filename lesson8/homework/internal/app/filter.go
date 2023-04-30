package app

import "time"

type Filter struct {
	AuthorID        int64
	PublishedAfter  time.Time
	PublishedBefore time.Time
}

type FilterOption func(*Filter)

func NewFilter(options ...FilterOption) *Filter {
	filter := &Filter{
		AuthorID:        -1,
		PublishedAfter:  time.Time{},
		PublishedBefore: time.Time{},
	}

	for _, option := range options {
		option(filter)
	}

	return filter
}

func WithAuthorID(author_id int64) FilterOption {
	return func(filter *Filter) {
		filter.AuthorID = author_id
	}
}

func WithPublishedAfter(publishedAfter time.Time) FilterOption {
	return func(filter *Filter) {
		filter.PublishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore time.Time) FilterOption {
	return func(filter *Filter) {
		filter.PublishedBefore = publishedBefore
	}
}
