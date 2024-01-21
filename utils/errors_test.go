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
		It("should panic", func() {
			Expect(func() { HandleError(errors.New("Should have panicked")) }).To(Panic())
		})
	})
})
