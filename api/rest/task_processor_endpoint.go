package rest

import (
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"danielpenchev98/http-job-processing-service/pkg/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Scheduler interface {
	CreateExecutionSequence(tasks []model.Task) ([]model.Task, error)
}

type ContentCreator interface {
	Create(tasks []model.Task) (string, error)
}

// Endpoint used for creation of execution plan and bash script content
type TaskProcessorEndpoint struct {
	scheduler Scheduler
	creator   ContentCreator
}

func NewTaskProcessorEndpoint(scheduler Scheduler, creator ContentCreator) *TaskProcessorEndpoint {
	return &TaskProcessorEndpoint{
		scheduler: scheduler,
		creator:   creator,
	}
}

// Creation of execution plan
func (e *TaskProcessorEndpoint) CreateProcessingPlan(c *gin.Context) {
	var content RequestPayload
	if err := c.ShouldBindJSON(&content); err != nil {
		sendErrorResponse(c, myerrors.NewClientError("Problem with the structure of the request body"))
		return
	}

	orderedTasks, err := e.createProcessingPlan(&content)
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	response := make([]Task, 0)
	for _, task := range orderedTasks {
		response = append(response, Task{
			Name:    task.Name,
			Command: task.Command,
		})
	}

	c.JSON(http.StatusOK, ExecutionPlanResponse{
		Status: http.StatusOK,
		Tasks:  response,
	})
}

// Creation of bash content 
func (e *TaskProcessorEndpoint) GenerateBashContent(c *gin.Context) {
	var content RequestPayload
	if err := c.ShouldBindJSON(&content); err != nil {
		sendErrorResponse(c, myerrors.NewClientError("Problem with the structure of the request body"))
		return
	}

	responseTasks, err := e.createProcessingPlan(&content)
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	bashContent, err := e.creator.Create(responseTasks)
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	c.String(http.StatusOK, bashContent)
}

func (e *TaskProcessorEndpoint) createProcessingPlan(requestContent *RequestPayload) ([]model.Task, error) {

	// Separate the request body representation from the inner representation of tasks
	tasks := make([]model.Task, 0)
	for _, task := range requestContent.Tasks {
		tasks = append(tasks, model.Task{
			Name:         task.Name,
			Command:      task.Command,
			Dependencies: task.DependableTasks,
		})
	}

	return e.scheduler.CreateExecutionSequence(tasks)
}

func sendErrorResponse(c *gin.Context, err error) {
	errorCode, errorMsg := getErrorResponseArguments(err)
	c.JSON(errorCode, ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	})
}

func getErrorResponseArguments(err error) (int, string) {
	var (
		errorCode int
		errorMsg  string
	)

	switch err.(type) {
	case *myerrors.ClientError:
		errorCode = http.StatusBadRequest
		errorMsg = fmt.Sprintf("Invalid request. Reason :%s", err.Error())
	default:
		log.Println(err)
		errorCode = http.StatusInternalServerError
		errorMsg = "Problem with the server, please try again later"
	}

	return errorCode, errorMsg
}
