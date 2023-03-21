package storage

import (
	"context"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	// TODO: implement this

	var (
		totalSize  int64 = 0
		totalCount int64 = 0
		curDir     Dir
	)

	queue := make([]Dir, 0)

	queue = append(queue, d)

	for len(queue) > 0 {

		curDir = queue[0]

		queue = queue[1:]

		dirs, files, err := curDir.Ls(ctx)

		if err != nil {
			return Result{}, err
		}

		queue = append(queue, dirs...)

		for _, file := range files {

			fileSize, err := file.Stat(ctx)

			if err != nil {
				return Result{}, err
			}

			totalSize += fileSize

			totalCount += 1

		}

	}

	return Result{Size: totalSize, Count: totalCount}, nil
}
