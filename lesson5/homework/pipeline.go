package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func h(ctx context.Context) Stage {
	return func(in In) Out {
		out := make(chan any)

		go func(out chan any) {
			defer close(out)
			for range ctx.Done() {
			}
		}(out)

		go func(in In, out chan any) {
			defer close(out)
			for {
				v, ok := <-in
				if !ok {
					return
				}
				out <- v
			}
		}(in, out)
		return out
	}
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {

	stages = append(stages, h(ctx))

	stagesRes := make([]Out, len(stages)+1)

	stagesRes[0] = in

	for i, st := range stages {
		stagesRes[i+1] = st(stagesRes[i])
	}

	return stagesRes[len(stagesRes)-1]

}
