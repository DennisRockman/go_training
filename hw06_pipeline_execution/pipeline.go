package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var doneChan = make(In)

func DoneChannel() In {
	return doneChan
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	doneChan = done

	for _, stage := range stages {
		in = stage(in)
	}

	return in // last result channel on runtime here
}
