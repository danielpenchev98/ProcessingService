package algorithm

import (
	"danielpenchev98/http-job-processing-service/pkg/model"
	"fmt"

	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
)

//go:generate mockgen --source=graph.go --destination algorithm_mocks/graph.go --package algorithm_mocks
type Graph interface {
	GetNodes() []string
	GetParents(string) ([]string, error)
}

// Used to represent the set of tasks as a graph
type DependencyGraph struct {
	nodes       []string
	parents map[string][]string
}

type DependencyGraphCreator struct {}

func (DependencyGraphCreator) Create(tasks []model.Task) (Graph, error) {
	if tasks == nil {
		return nil, myerrors.NewServerError("invalid task sequence supplied to dependency graph factory")
	}

	nodes := make([]string, 0)
	parents := make(map[string][]string)

	visited := make(map[string]struct{})

	for _, task := range tasks {
		nodes = append(nodes, task.Name)
		edges := make([]string, 0)

		if _, ok := visited[task.Name]; ok {
			return nil, myerrors.NewClientError("multiple tasks with the same name")
		}

		visited[task.Name] = struct{}{}

		if task.Dependencies == nil {
			parents[task.Name] = edges
			continue
		}

		edges = append(edges, task.Dependencies...)
		parents[task.Name] = edges
	}

	return &DependencyGraph{
		nodes:       nodes,
		parents: parents,
	}, nil
}

// Getting dependencies of a particular task
func (g *DependencyGraph) GetParents(nodeName string) ([]string, error) {
	if _, ok := g.parents[nodeName]; !ok {
		return nil, myerrors.NewClientError(fmt.Sprintf("dependency [%s] not found", nodeName))
	}
	return g.parents[nodeName], nil
}

// Getting all nodes from the graph
func (g *DependencyGraph) GetNodes() []string {
	var result []string
	return append(result, g.nodes...)
}
