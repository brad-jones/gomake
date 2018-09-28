package generator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getDefaultValueForType", func() {

	When("given an unrecognised type", func() {
		It("should return nil", func() {
			Expect(getDefaultValueForType("foo", false)).To(Equal("nil"))
		})
	})

	When("given an array type", func() {
		It("should return nil", func() {
			Expect(getDefaultValueForType("foo", true)).To(Equal("nil"))
		})
	})

	When("given a string type", func() {
		It("should return an empty string", func() {
			Expect(getDefaultValueForType("string", false)).To(Equal(`""`))
		})
	})

	When("given a bool type", func() {
		It("should return false", func() {
			Expect(getDefaultValueForType("bool", false)).To(Equal("false"))
		})
	})

	When("given a number type", func() {
		It("should return zero", func() {
			Expect(getDefaultValueForType("int", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int16", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint8", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint16", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("float32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("float64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("Duration", false)).To(Equal("0"))
		})
	})

})
