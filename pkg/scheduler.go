package pkg

import (
	"danielpenchev98/http-job-processing-service/pkg/algorithm"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"danielpenchev98/http-job-processing-service/pkg/model"
)

//go:generate mockgen --source=scheduler.go --destination pkg_mocks/scheduler.go --package pkg_mocks
type GraphAlgorithm interface {
	Apply(algorithm.Graph) ([]string, error)
}

type GraphCreator interface {
	Create(task []model.Task) (algorithm.Graph, error)
}

type TaskScheduler struct {
	algorithm    GraphAlgorithm
	graphCreator GraphCreator
}

func NewTaskScheduler(algorithm GraphAlgorithm, graphCreator GraphCreator) *TaskScheduler {
	return &TaskScheduler{
		algorithm:    algorithm,
		graphCreator: graphCreator,
	}
}

func (s *TaskScheduler) CreateExecutionSequence(tasks []model.Task) ([]model.Task, error) {
	graph, err := s.graphCreator.Create(tasks)
	if err != nil {
		return nil, handleError(err)
	}

	order, err := s.algorithm.Apply(graph)
	if err != nil {
		return nil, handleError(err)
	}

	mapping := make(map[string]model.Task)
	for _, task := range tasks {
		mapping[task.Name] = task
	}

	result := make([]model.Task, 0)
	for _, taskName := range order {
		result = append(result, mapping[taskName])
	}

	return result, nil
}

func handleError(err error) error {
	msg := "problem with creation of execution sequence"
	switch err.(type) {
	case *myerrors.ClientError:
		return myerrors.NewClientErrorWrap(err, msg)
	default:
		return myerrors.NewServerErrorWrap(err, msg)
	}
}
