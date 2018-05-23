package routine

import (
	"sync"
)

type Job struct {
	main func()
}

func NewJob(main func()) *Job {
	return &Job{main}
}

func NewParallelJob(parallel int, main func()) *Job {
	return &Job{func() {
		wg := &sync.WaitGroup{}
		wg.Add(parallel)
		for i := 0; i < parallel; i++ {
			go func() {
				defer wg.Done()
				main()
			}()
		}
		wg.Wait()
	}}
}

func (j *Job) Then(main func()) *Job {
	var last = j.main
	return &Job{func() { last(); main() }}
}

func (j *Job) Run() <-chan struct{} {
	var done = make(chan struct{})
	go func() {
		defer close(done)
		j.main()
	}()
	return done
}

func (j *Job) RunAndWait() {
	<-j.Run()
}
