package model

// Representation of a Task, used for computation of tasks
type Task struct {
	Name         string
	Command      string
	Dependencies []string
}
