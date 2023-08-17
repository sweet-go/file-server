package model

// Task defined any task to be registered to worker
type Task string

// list of tasks
var (
	TaskScheduledDelete Task = "SCHEDULED_DELETE"
)
