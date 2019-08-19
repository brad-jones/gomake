package executor_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExecutor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestExecutor")
}

/*
var _ = Describe("findGomakeFolder", func() {

	When("given a dir that contains a .gomake folder", func() {
		It("should return the absolute path to that folder", func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			result, err := findGomakeFolder(filepath.Join(cwd, "..", "example"))
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(filepath.Join(cwd, "..", "example", ".gomake")))
		})
	})

	When("given a dir that contains a .gomake folder in a parent dir", func() {
		It("should return the absolute path to that folder", func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			result, err := findGomakeFolder(filepath.Join(cwd, "..", "example", "a-project", "src"))
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(filepath.Join(cwd, "..", "example", ".gomake")))
		})
	})

	When("given a dir that does not contain a .gomake folder, nor any parent dir", func() {
		It("should return a ErrReachedRootOfFs error", func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			_, err = findGomakeFolder(filepath.Join(cwd))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: failed to find valid '.gomake' folder, reached root of filesystem"))
		})
	})

})
*/
