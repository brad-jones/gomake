package generator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mapTypeToFlagType", func() {

	When("given bool", func() {
		It("should return BoolP", func() {
			result, err := mapTypeToFlagType("bool", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BoolP"))
		})
	})

	When("given bool and is array", func() {
		It("should return BoolSliceP", func() {
			result, err := mapTypeToFlagType("bool", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BoolSliceP"))
		})
	})

	When("given string", func() {
		It("should return StringP", func() {
			result, err := mapTypeToFlagType("string", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("StringP"))
		})
	})

	When("given string and is array", func() {
		It("should return StringSliceP", func() {
			result, err := mapTypeToFlagType("string", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("StringSliceP"))
		})
	})

	When("given int", func() {
		It("should return IntP", func() {
			result, err := mapTypeToFlagType("int", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IntP"))
		})
	})

	When("given int and is array", func() {
		It("should return IntSliceP", func() {
			result, err := mapTypeToFlagType("int", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IntSliceP"))
		})
	})

	When("given int8", func() {
		It("should return Int8P", func() {
			result, err := mapTypeToFlagType("int8", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int8P"))
		})
	})

	When("given int16", func() {
		It("should return Int16P", func() {
			result, err := mapTypeToFlagType("int16", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int16P"))
		})
	})

	When("given int32", func() {
		It("should return Int32P", func() {
			result, err := mapTypeToFlagType("int32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int32P"))
		})
	})

	When("given int64", func() {
		It("should return Int64P", func() {
			result, err := mapTypeToFlagType("int64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int64P"))
		})
	})

	When("given float32", func() {
		It("should return Float32P", func() {
			result, err := mapTypeToFlagType("float32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Float32P"))
		})
	})

	When("given float64", func() {
		It("should return Float64P", func() {
			result, err := mapTypeToFlagType("float64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Float64P"))
		})
	})

	When("given uint", func() {
		It("should return UintP", func() {
			result, err := mapTypeToFlagType("uint", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("UintP"))
		})
	})

	When("given uint and is array", func() {
		It("should return UintSliceP", func() {
			result, err := mapTypeToFlagType("uint", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("UintSliceP"))
		})
	})

	When("given uint8", func() {
		It("should return Uint8P", func() {
			result, err := mapTypeToFlagType("uint8", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint8P"))
		})
	})

	When("given uint16", func() {
		It("should return Uint16P", func() {
			result, err := mapTypeToFlagType("uint16", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint16P"))
		})
	})

	When("given uint32", func() {
		It("should return Uint32P", func() {
			result, err := mapTypeToFlagType("uint32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint32P"))
		})
	})

	When("given uint64", func() {
		It("should return Uint64P", func() {
			result, err := mapTypeToFlagType("uint64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint64P"))
		})
	})

	When("given Duration", func() {
		It("should return DurationP", func() {
			result, err := mapTypeToFlagType("Duration", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("DurationP"))
		})
	})

	When("given Duration and is array", func() {
		It("should return DurationSliceP", func() {
			result, err := mapTypeToFlagType("Duration", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("DurationSliceP"))
		})
	})

	When("given IP", func() {
		It("should return IPP", func() {
			result, err := mapTypeToFlagType("IP", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPP"))
		})
	})

	When("given IP and is array", func() {
		It("should return IPSliceP", func() {
			result, err := mapTypeToFlagType("IP", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPSliceP"))
		})
	})

	When("given IPMask", func() {
		It("should return IPMaskP", func() {
			result, err := mapTypeToFlagType("IPMask", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPMaskP"))
		})
	})

	When("given byte and is array", func() {
		It("should return BytesBase64P", func() {
			result, err := mapTypeToFlagType("byte", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BytesBase64P"))
		})
	})

	When("given foo an unsupported type", func() {
		It("should return ErrUnsupportedFlagType", func() {
			_, err := mapTypeToFlagType("foo", false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: the flag type 'FooP' is unsupported by gomake and/or cobra"))
		})
	})

})
