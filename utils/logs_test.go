package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Local module
	. "github.com/fleroy-isagri/env2js/utils"
)

var _ = Describe("Logs", func() {
	It("should not panic when calling the LogError function", func() {
		Expect(func() { LogError("Test", "Error") }).NotTo(Panic())
	})

	It("should not panic when calling the LogSucess function", func() {
		Expect(func() { LogSuccess("Test : ", "Sucess") }).NotTo(Panic())
	})
})
