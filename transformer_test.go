package main

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"gopkg.in/yaml.v2"
)

func TestTransformer(t *testing.T) {
	spec.Run(t, "Transformer", func(t *testing.T, when spec.G, it spec.S) {
		var input map[interface{}]interface{}
		var output BulkImport

		it.Before(func() {
			RegisterTestingT(t)
			file, _ := os.Open("fixtures/input.yml")
			inputDecoder := yaml.NewDecoder(file)
			inputDecoder.Decode(&input)

			file, _ = os.Open("fixtures/output.yml")

			outputDecoder := yaml.NewDecoder(file)
			outputDecoder.Decode(&output)
		})

		when("Doing an ETL", func() {
			it("Imported the input yml successfully", func() {
				Expect(len(input)).To(Equal(8))
			})

			it("Imported the output yml successfully", func() {
				Expect(len(output.Credentials)).To(Equal(8))
			})

			it("Transformed the output correctly", func() {
				file, _ := os.Open("fixtures/input.yml")
				out, err := Transform("/foo", file)

				Expect(err).To(BeNil())
				Expect(len(out.Credentials)).To(Equal(8))
			})

			it("Failed to transform bad input", func() {
				file, _ := os.Open("fixtures/input_invalid.yml")
				out, err := Transform("/foo", file)
				Expect(err).To(Not(BeNil()))
				Expect(len(out.Credentials)).To(Equal(0))
			})
		})
	})
}
