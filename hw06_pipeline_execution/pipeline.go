package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		tempCh := make(Bi)

		go func(in In) {
			// Слушаем in, отправляем в tempCh
		}(in)

		in = stage(tempCh)
	}

	return in
}
