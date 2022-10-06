package script_test

import (
	"danielpenchev98/http-job-processing-service/pkg/script"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const scriptStartInstruction = "#!/usr/bin/env bash\n\n"

var _ = Describe("BashContentBuilderCreator", func() {
	Context("Create", func() {
		When("creating a new instance of bash content builder", func(){
			It("succeeds", func(){
				builderCreator := script.BashContentBuilderCreator{}
				builder := builderCreator.Create()
				Expect(builder).NotTo(BeNil())
			})
		})
	})
})

var _ = Describe("ScriptContentBuilder", func() {
	var builder script.ContentBuilder

	BeforeEach(func() {
		builderCreator := script.BashContentBuilderCreator{}
		builder = builderCreator.Create()
	})

	When("commands are added", func() {
		It("contains all commands", func() {
			builder.AddCommand("cmd1").
				AddCommand("cmd2")

			content, err := builder.String()
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(createScriptContent("cmd1", "cmd2")))
		})
	})

	When("builder is created", func() {
		It("contains only the start instruction", func() {
			content, err := builder.String()
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(createScriptContent()))
		})
	})
})

func createScriptContent(cmds ...string) string {
	content := scriptStartInstruction
	for _, cmd := range cmds {
		content += cmd + "\n"
	}
	return content
}
