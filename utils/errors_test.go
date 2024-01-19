package utils_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Local module
	. "github.com/fleroy-isagri/env2js/utils"
)

var _ = Describe("Errors", func() {
	Context("when calling the HandleError function", func() {
		doesPanic := true
		It("should panic", func() {
			defer func() {
				recover()
				Expect(doesPanic).To(Equal(true))
			}()
			HandleError(errors.New("Should have panicked"))
			doesPanic = false
		})
	})
})
