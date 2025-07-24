package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		tempCh := make(Bi)

		go func(ch In) {
			defer func() {
				go func() {
					for {
						_, ok := <-ch

						if !ok {
							return
						}
					}
				}()
				close(tempCh)
			}()

			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						return
					}
					tempCh <- msg
				case <-done:
					return
				}
			}
		}(out)

		out = stage(tempCh)
	}

	return out
}
