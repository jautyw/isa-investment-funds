package test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func InvestmentsTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Investments Suite")
}
