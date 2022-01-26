package scheduler

type Scheduler interface {
	Schedule(fn func()) error
}
