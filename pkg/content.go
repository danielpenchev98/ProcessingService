package pkg

import (
	"danielpenchev98/http-job-processing-service/pkg/model"
	"danielpenchev98/http-job-processing-service/pkg/script"
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
)

//go:generate mockgen --source=content.go --destination pkg_mocks/content.go --package pkg_mocks
type ContentBuilderCreator interface {
	Create() script.ContentBuilder
}

type ContentCreator struct {
	builderCreator ContentBuilderCreator
}

func NewContentCreator(contentBuilderFactory ContentBuilderCreator) *ContentCreator {
	return &ContentCreator{
		builderCreator: contentBuilderFactory,
	}
}

func (c *ContentCreator) Create(tasks []model.Task) (string, error) {
	if tasks == nil {
		return "", myerrors.NewServerError("invalid tasks array supplied")
	}

	builder := c.builderCreator.Create()

	for _, task := range tasks {
		builder.AddCommand(task.Command)
	}

	content, err := builder.String()
	if err != nil {
		return "",  myerrors.NewServerErrorWrap(err, "problem with content creation")
	}
	return content, nil
}
