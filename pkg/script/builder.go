package script

import (
	myerrors "danielpenchev98/http-job-processing-service/pkg/errors"
	"strings"
)

const scriptStartInstruction = "#!/usr/bin/env bash\n"

//go:generate mockgen --source=builder.go --destination script_mocks/builder.go --package script_mocks
type ContentBuilder interface {
	AddCommand(string) ContentBuilder
	String() (string, error)
}


// Builder pattern for bash script content creation
type BashContentBuilder struct {
	content []string
}

type BashContentBuilderCreator struct{}

func (BashContentBuilderCreator) Create() ContentBuilder {
	return &BashContentBuilder{
		content: []string{scriptStartInstruction},
	}
}

// Adds the next command in the script
func (b *BashContentBuilder) AddCommand(cmd string) ContentBuilder {
	b.content = append(b.content, cmd)
	return b
}

// Builds the content of the bash script
func (b *BashContentBuilder) String() (string, error) {
	builder := strings.Builder{}
	for _, cmd := range b.content {
		if _, err := builder.WriteString(cmd + "\n"); err != nil {
			return "", myerrors.NewServerErrorWrap(err, "problem with concatenating all commands/instructions into a single string")
		}
	}

	return builder.String(), nil
}
