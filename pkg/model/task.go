package model

type Task struct {
	Name         string
	Command      string
	Dependencies []string
}
