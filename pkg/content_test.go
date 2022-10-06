package pkg_test

import (
	"danielpenchev98/http-job-processing-service/pkg"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"danielpenchev98/http-job-processing-service/pkg/model"
	"danielpenchev98/http-job-processing-service/pkg/pkg_mocks"
	"danielpenchev98/http-job-processing-service/pkg/script/script_mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContentCreator", func() {

	var (
		contentCreator *pkg.ContentCreator
		builderCreator *pkg_mocks.MockContentBuilderCreator
		builder *script_mocks.MockContentBuilder
	)

	BeforeEach(func() {
		controller := gomock.NewController(GinkgoT())
		builderCreator = pkg_mocks.NewMockContentBuilderCreator(controller)
		builder = script_mocks.NewMockContentBuilder(controller)
	})

	Context("NewContentCreator", func() {
		When("creating a new instance of content creator", func() {
			It("returns a content creator instance", func() {
				contentCreator = pkg.NewContentCreator(builderCreator)
				Expect(contentCreator).NotTo(BeNil())
			})
		})
	})

	Context("Create", func() {

		BeforeEach(func() {
			contentCreator = pkg.NewContentCreator(builderCreator)
		})

		When("nil tasks are supplied", func() {
			It("returns error", func() {
				_, err := contentCreator.Create(nil)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ServerError)
				Expect(ok).To(BeTrue())

			})
		})

		When("content building failed", func() {
			It("returns error", func() {
				tasks := []model.Task{{Name: "task1",Command: "cmd1"},{Name: "task2",Command: "cmd2"}}

				builderCreator.EXPECT().Create().Return(builder)
				builder.EXPECT().AddCommand("cmd1").Return(builder)
				builder.EXPECT().AddCommand("cmd2").Return(builder)
				builder.EXPECT().String().Return("", myerrors.NewServerError(""))

				_, err := contentCreator.Create(tasks)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerrors.ServerError)
				Expect(ok).To(BeTrue())
			})
		})

		When("content building succeeded", func() {
			It("succeeds", func() {
				tasks := []model.Task{{Name: "task1",Command: "cmd1"},{Name: "task2",Command: "cmd2"}}

				builderCreator.EXPECT().Create().Return(builder)
				builder.EXPECT().AddCommand("cmd1").Return(builder)
				builder.EXPECT().AddCommand("cmd2").Return(builder)

				expectedContent := "expectedContent"
				builder.EXPECT().String().Return(expectedContent, nil)

				realContent, err := contentCreator.Create(tasks)
				Expect(err).NotTo(HaveOccurred())
				Expect(realContent).To(Equal(expectedContent))
			})
		})
	})
})
