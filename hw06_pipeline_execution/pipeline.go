package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		in = link(in, done)
	}

	for _, stage := range stages {
		if stage != nil {
			in = stage(link(in, done))
		}
	}
	return in
}

func link(in In, done In) Bi {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
