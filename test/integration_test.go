package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var port string

var _ = BeforeSuite(func() {
	const defaultPort = "8080"

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//Workaround though bash to start the application
	cmd := exec.Command("./setup_integration_test.sh")
	_, err := cmd.Output()
	if err != nil {
		msg := fmt.Sprintf("problem with integration test setup - reason: [%v]", err)
		Fail(msg)
	}
})

var _ = AfterSuite(func() {
	cmd := exec.Command("./clean_integration_test.sh")
	_, err := cmd.Output()
	if err != nil {
		msg := fmt.Sprintf("problem with integration test cleanup - reason: [%v]", err)
		Fail(msg)
	}
})

var _ = Describe("HealthCheckEndpoint", func() {

	When("service is available", func() {
		It("returns success", func() {
			resp, err := http.Get(getHealthChecktURL())
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("TaskProcessorEndpoint", func() {

	const (
		jsonContentType = "application/json"
		requestDirPath  = "request"
		responseDirPath = "response"
	)

	var (
		reqPayload          []byte
		expectedRespPayload []byte
	)

	Context("/v1/processing-plan", func() {
		When("there is a circular dependency between tasks", func() {

			BeforeEach(func() {
				fileName := "circular_dependency.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getProcessingPlanURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("exists a dependency to non existing task", func() {

			BeforeEach(func() {
				fileName := "non_existing_dependency.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getProcessingPlanURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("invalid type of payload is used", func() {
			BeforeEach(func() {
				reqPayload = readResourceContent(requestDirPath + "/invalid_structure.xml")
				expectedRespPayload = readResourceContent(responseDirPath + "/invalid_structure.json")
			})

			It("returns error", func() {
				resp, err := http.Post(getProcessingPlanURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("multiple tasks with the same name are supplied", func() {
			BeforeEach(func() {
				fileName := "same_task_names.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getProcessingPlanURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("valid payload is supplied", func() {
			BeforeEach(func() {
				fileName := "valid_set_tasks.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("succeeds", func() {
				resp, err := http.Post(getProcessingPlanURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

	})

	Context("/v1/bash", func() {
		When("there is a circular dependency between tasks", func() {

			BeforeEach(func() {
				fileName := "circular_dependency.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getBashURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("exists a dependency to non existing task", func() {

			BeforeEach(func() {
				fileName := "non_existing_dependency.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getBashURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("invalid type of payload is used", func() {
			BeforeEach(func() {
				reqPayload = readResourceContent(requestDirPath + "/invalid_structure.xml")
				expectedRespPayload = readResourceContent(responseDirPath + "/invalid_structure.json")
			})

			It("returns error", func() {
				resp, err := http.Post(getBashURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("multiple tasks with the same name are supplied", func() {
			BeforeEach(func() {
				fileName := "same_task_names.json"
				reqPayload = readResourceContent(requestDirPath + "/" + fileName)
				expectedRespPayload = readResourceContent(responseDirPath + "/" + fileName)
			})

			It("returns error", func() {
				resp, err := http.Post(getBashURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})

		When("valid payload is supplied", func() {
			BeforeEach(func() {
				reqPayload = readResourceContent(requestDirPath + "/valid_set_tasks.json")
				expectedRespPayload = readResourceContent(responseDirPath + "/valid_set_tasks.sh")
			})

			It("succeeds", func() {
				resp, err := http.Post(getBashURL(), jsonContentType, bytes.NewBuffer(reqPayload))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				content, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(Equal(expectedRespPayload))
			})
		})
	})
})

func getHealthChecktURL() string {
	return fmt.Sprintf("http://localhost:%s/v1/healthcheck", port)
}

func getProcessingPlanURL() string {
	return fmt.Sprintf("http://localhost:%s/v1/processing-plan", port)
}

func getBashURL() string {
	return fmt.Sprintf("http://localhost:%s/v1/bash", port)
}

func readResourceContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Fail(fmt.Sprintf("Problem with reading file content - reason: [%v]", err))
	}
	return content
}
