package algorithm

import (
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"fmt"
)

const (
	halfVisited  = 0 // some of its neightbours aren't fully visited and the current node is visited
	fullyVisited = 1 // all of his neightbours are fully visited and the current node is visited
)

// Algorithm for sorting graph nodes such that every node V in the sequence is dependent only on some of the previous nodes in the sequence 
type TopologicalSort struct{}

func (TopologicalSort) Apply(taskGraph Graph) ([]string, error) {
	visited := make(map[string]int)
	result := make([]string, 0)

	for _, task := range taskGraph.GetNodes() {
		subResult, err := dfs(task, taskGraph, visited)
		if err != nil {
			return nil, err
		}
		result = append(result, subResult...)
	}

	return result, nil
}

func dfs(taskName string, graph Graph, visited map[string]int) ([]string, error) {
	if state, ok := visited[taskName]; ok {
		if state == halfVisited {
			return nil, myerrors.NewClientError("cyclic dependency detected")
		}
		return []string{}, nil
	}

	visited[taskName] = halfVisited
	result := make([]string, 0)

	tasks, err := graph.GetParents(taskName)
	if err != nil {
		return nil, myerrors.NewClientErrorWrap(err, fmt.Sprintf("problem with getting the dependencies of node [%s]", taskName))
	}

	for _, task := range tasks {
		childResult, err := dfs(task, graph, visited)
		if err != nil {
			return nil, err
		}

		result = append(result, childResult...)
	}

	result = append(result, taskName)
	visited[taskName] = fullyVisited
	return result, nil
}
