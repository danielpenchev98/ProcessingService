package pkg_test

import (
	"danielpenchev98/http-job-processing-service/pkg"
	"danielpenchev98/http-job-processing-service/pkg/algorithm/algorithm_mocks"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"danielpenchev98/http-job-processing-service/pkg/model"
	"danielpenchev98/http-job-processing-service/pkg/pkg_mocks"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskScheduler", func() {
	var (
		taskNames          []string
		tasks              []model.Task
		mockGraph          *algorithm_mocks.MockGraph
		mockGraphCreator   *pkg_mocks.MockGraphCreator
		mockGraphAlgorithm *pkg_mocks.MockGraphAlgorithm
		scheduler          *pkg.TaskScheduler
	)

	BeforeEach(func() {
		taskNames = []string{"task1", "task2", "task3"}
		tasks = []model.Task{{Name: taskNames[0]}, {Name: taskNames[1]},{Name: taskNames[2]}}
		controller := gomock.NewController(GinkgoT())
		mockGraph = algorithm_mocks.NewMockGraph(controller)

		mockGraphAlgorithm = pkg_mocks.NewMockGraphAlgorithm(controller)
		mockGraphCreator = pkg_mocks.NewMockGraphCreator(controller)
		scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
	})

	Context("CreateExecutionSequence", func() {
		When("graph creation failed", func() {
			Context("because of the client input", func() {
				It("propagates error", func() {
					mockGraphCreator.EXPECT().Create(tasks).Return(nil, myerrors.NewClientError(""))

					scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
					_, err := scheduler.CreateExecutionSequence(tasks)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerrors.ClientError)
					Expect(ok).To(BeTrue())
				})
			})

			Context("because of some system error", func() {
				It("propagates error", func() {
					mockGraphCreator.EXPECT().Create(tasks).Return(nil, errors.New(""))

					scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
					_, err := scheduler.CreateExecutionSequence(tasks)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerrors.ServerError)
					Expect(ok).To(BeTrue())
				})
			})
		})

		When("graph creation succeeded", func() {
			Context("and graph algorithm failed", func() {
				Context("because of client input error", func() {
					It("propagates error", func() {
						mockGraphCreator.EXPECT().Create(tasks).Return(mockGraph, nil)
						mockGraphAlgorithm.EXPECT().Apply(mockGraph).Return(nil, myerrors.NewClientError(""))

						scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
						_, err := scheduler.CreateExecutionSequence(tasks)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerrors.ClientError)
						Expect(ok).To(BeTrue())
					})
				})

				Context("because of some system error", func() {
					It("propagates error", func() {
						mockGraphCreator.EXPECT().Create(tasks).Return(mockGraph, nil)
						mockGraphAlgorithm.EXPECT().Apply(mockGraph).Return(nil, errors.New(""))

						scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
						_, err := scheduler.CreateExecutionSequence(tasks)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerrors.ServerError)
						Expect(ok).To(BeTrue())
					})
				})
			})

			Context("and graph algorithm succeeded", func() {
				It("returns ordered tasks", func() {
					order := []string{taskNames[2], taskNames[0], taskNames[1]}

					mockGraphCreator.EXPECT().Create(tasks).Return(mockGraph, nil)
					mockGraphAlgorithm.EXPECT().Apply(mockGraph).Return(order, nil)

					scheduler = pkg.NewTaskScheduler(mockGraphAlgorithm, mockGraphCreator)
					orderedTasks, err := scheduler.CreateExecutionSequence(tasks)
					Expect(err).NotTo(HaveOccurred())

					expectedOrderedTasks := []model.Task{tasks[2], tasks[0], tasks[1]}
					Expect(orderedTasks).To(Equal(expectedOrderedTasks))
				})
			})
		})
	})

})
