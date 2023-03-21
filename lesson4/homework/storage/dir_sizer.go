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

type dirInfo struct {
	totalFileSize  int64
	totalFileCount int64
	subDirsNum     int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func extract(ctx context.Context, d Dir, dirInfoCh chan dirInfo, errCh chan error) {

	ctx.Value(struct{}{}).(chan struct{}) <- struct{}{}

	dirs, files, err := d.Ls(ctx)

	if err != nil {
		errCh <- err
		return
	}

	var (
		totalFileSize  int64 = 0
		totalFileCount int64 = 0
		subDirsNum     int   = 0
	)

	for _, file := range files {

		fileSize, err := file.Stat(ctx)

		if err != nil {
			errCh <- err
			return
		}

		totalFileSize += fileSize
		totalFileCount++

	}

	for _, dir := range dirs {

		subDirsNum++
		go extract(ctx, dir, dirInfoCh, errCh)

	}

	dirInfoCh <- dirInfo{totalFileSize, totalFileCount, subDirsNum}

	<-ctx.Value(struct{}{}).(chan struct{})

}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	// TODO: implement this

	var (
		size       int64 = 0
		count      int64 = 0
		toExctract int   = 1
	)

	dirInfoCh := make(chan dirInfo)
	errCh := make(chan error)
	goroutineCh := make(chan struct{}, max(a.maxWorkersCount, 4))
	ctx = context.WithValue(ctx, struct{}{}, goroutineCh)

	go extract(ctx, d, dirInfoCh, errCh)
	for n := 0; n < toExctract; n++ {
		select {
		case d := <-dirInfoCh:
			size += d.totalFileSize
			count += d.totalFileCount
			toExctract += d.subDirsNum
		case e := <-errCh:
			return Result{}, e
		}
	}

	return Result{Size: size, Count: count}, nil
}
