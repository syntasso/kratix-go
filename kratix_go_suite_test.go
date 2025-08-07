package kratix_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKratixGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KratixGo Suite")
}
