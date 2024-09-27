package gojm

import priorityqueue "github.com/xybor-x/priority_queue"

type JobWrapper struct {
	Priority         Priority
	OriginalPriority Priority
	job              *Job
}

func wrapJob(e priorityqueue.Element[*Job]) JobWrapper {
	return JobWrapper{
		OriginalPriority: e.OriginalPriority().(Priority),
		Priority:         e.Priority().(Priority),
		job:              e.To(),
	}
}

func (w JobWrapper) Unwrap() *Job {
	return w.job
}
