package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEnv2js(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils suite")
}
