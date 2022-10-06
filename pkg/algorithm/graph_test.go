package algorithm_test

import (
	"danielpenchev98/http-job-processing-service/pkg/algorithm"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"danielpenchev98/http-job-processing-service/pkg/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var taskNames = []string{"task1", "task2", "task3", "task4"}

var _ = Describe("DependencyGraphCreator", func() {

	var creator algorithm.DependencyGraphCreator

	Context("Create", func() {
		When("multiple nodes with the same name exist", func() {
			It("should return error", func() {
				tasks := make([]model.Task, 0)
				tasks = append(tasks, model.Task{Name: taskNames[0]})
				tasks = append(tasks, model.Task{Name: taskNames[0]})

				_, err := creator.Create(tasks)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ClientError)
				Expect(ok).To(Equal(true))

			})
		})

		When("nil task list is provided", func() {
			It("should return error", func() {
				_, err := creator.Create(nil)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("non nil task list is provided", func() {
			It("should succeed", func() {
				tasks := make([]model.Task, 0)
				tasks = append(tasks, model.Task{Name: taskNames[0]})

				graph, err := creator.Create(tasks)
				Expect(err).NotTo(HaveOccurred())
				Expect(graph).NotTo(BeNil())
			})
		})
	})
})

var _ = Describe("DependencyGraph", func() {
	var (
		graph        algorithm.Graph
		graphCreator algorithm.DependencyGraphCreator
	)

	Context("GetNeighbours", func() {

		BeforeEach(func() {
			tasks := make([]model.Task, 0)
			tasks = append(tasks, model.Task{
				Name:         taskNames[0],
				Dependencies: []string{taskNames[1], taskNames[2]},
			})

			tasks = append(tasks, model.Task{
				Name:         taskNames[1],
				Dependencies: []string{taskNames[3]},
			})

			tasks = append(tasks, model.Task{
				Name:         taskNames[2],
				Dependencies: []string{taskNames[3]},
			})

			tasks = append(tasks, model.Task{
				Name:         taskNames[3],
				Dependencies: []string{},
			})

			var err error
			graph, err = graphCreator.Create(tasks)
			Expect(err).NotTo(HaveOccurred())
		})

		When("node isnt in the graph", func() {
			It("should return error", func() {
				_, err := graph.GetParents("invalid node name")
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ClientError)
				Expect(ok).To(Equal(true))
			})
		})

		When("node is in the graph", func() {
			Context("and has no neighbours", func() {
				It("should return empty list", func() {
					neighbours, err := graph.GetParents(taskNames[3])
					Expect(err).NotTo(HaveOccurred())
					Expect(neighbours).To(BeEmpty())
				})
			})

			Context("and has neighbours", func() {
				It("should return all of his neighbours", func() {
					neighbours, err := graph.GetParents(taskNames[0])
					Expect(err).NotTo(HaveOccurred())
					Expect(neighbours).To(ContainElements(taskNames[1], taskNames[2]))
				})
			})
		})
	})

	Context("GetNodes", func() {
		When("graph has not nodes", func() {
			It("should return empty list", func() {
				graph, err := graphCreator.Create([]model.Task{})
				Expect(err).NotTo(HaveOccurred())
				nodes := graph.GetNodes()
				Expect(nodes).To(BeEmpty())
			})
		})

		When("graph has nodes", func() {
			It("should return all graph nodes", func() {
				tasks := make([]model.Task, 0)
				tasks = append(tasks, model.Task{
					Name: taskNames[0],
				})

				tasks = append(tasks, model.Task{
					Name: taskNames[1],
				})

				graph, _ = graphCreator.Create(tasks)

				nodes := graph.GetNodes()
				Expect(nodes).To(ContainElements(taskNames[0], taskNames[1]))
			})
		})
	})
})
