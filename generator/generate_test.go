package generator_test

import (
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/brad-jones/gomake.v2/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {

	When("given an example gomake dir", func() {
		It("should generate a runable cobra cli app", func() {
			cwd, _ := os.Getwd()
			dir := filepath.Join(cwd, "..", "example", ".gomake")
			err := generator.Generate(dir)
			Expect(err).ToNot(HaveOccurred())

			out, err := exec.Command("go", "run", dir).Output()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(ContainSubstring("Makefile written in golang"))
		})
	})

})
