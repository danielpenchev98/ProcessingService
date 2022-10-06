package algorithm_test

import (
	"danielpenchev98/http-job-processing-service/pkg/algorithm"
	"danielpenchev98/http-job-processing-service/pkg/algorithm/algorithm_mocks"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DependencyGraph", func() {
	var (
		graph *algorithm_mocks.MockGraph
		graphAlgorithm algorithm.TopologicalSort
	)

	BeforeEach(func() {
		controller := gomock.NewController(GinkgoT())
		graph = algorithm_mocks.NewMockGraph(controller)
		graphAlgorithm = algorithm.TopologicalSort{}
	})

	Context("TopologicalSort", func() {
		When("graph is empty", func() {
			It("succeeds", func() {
				graph.EXPECT().GetNodes().Return([]string{})
				result, err := graphAlgorithm.Apply(graph)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(BeEmpty())
			})
		})

		When("graph has a cycle", func() {
			It("returns an error", func() {
				graph.EXPECT().GetNodes().Return([]string{"task1", "task2"})
				graph.EXPECT().GetParents("task1").Return([]string{"task2"}, nil)
				graph.EXPECT().GetParents("task2").Return([]string{"task1"}, nil)

				_, err := graphAlgorithm.Apply(graph)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ClientError)
				Expect(ok).To(Equal(true))
			})
		})

		When("graph node is dependent on a node outside the graph", func() {
			It("returns an error", func() {
				graph.EXPECT().GetNodes().Return([]string{"task1"})
				graph.EXPECT().GetParents("task1").Return([]string{"task2"}, nil)
				graph.EXPECT().GetParents("task2").Return(nil, myerrors.NewClientError(""))

				_, err := graphAlgorithm.Apply(graph)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ClientError)
				Expect(ok).To(BeTrue())
			})
		})

		When("graph is acyclic", func() {
			It("succeeds", func() {
				graph.EXPECT().GetNodes().Return([]string{"task1", "task2", "task3", "task4"})
				graph.EXPECT().GetParents("task1").Return([]string{}, nil)
				graph.EXPECT().GetParents("task2").Return([]string{"task1"}, nil)
				graph.EXPECT().GetParents("task3").Return([]string{"task1"}, nil)
				graph.EXPECT().GetParents("task4").Return([]string{"task2", "task3"}, nil)

				result, err := graphAlgorithm.Apply(graph)
				Expect(err).NotTo(HaveOccurred())

				topologicalOrder := []string{"task1", "task2", "task3", "task4"}
				Expect(result).To(Equal(topologicalOrder))
			})
		})
	})

})
