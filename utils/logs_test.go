package utils_test

import (
	. "github.com/onsi/ginkgo/v2"

	// Local module
	. "github.com/fleroy-isagri/env2js/utils"
)

var _ = Describe("Logs", func() {
	Context("when calling the LogError function", func() {
		It("should not panic", func() {
			defer CheckIfPanic("LogError panic")
			LogError("Test", "Error")
		})
	})

	Context("when calling the LogSucess function", func() {
		It("should not panic", func() {
			defer CheckIfPanic("LogSuccess panic")
			LogSuccess("Test : ", "Sucess")
		})
	})
})

func CheckIfPanic(panicMessage string) {
	if r := recover(); r != nil {
		AbortSuite(panicMessage)
	}
}
