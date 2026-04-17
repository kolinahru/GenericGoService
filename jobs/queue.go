package jobs

type Queue struct {
	Jobs chan Job
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{
		Jobs: make(chan Job, bufferSize),
	}
}
