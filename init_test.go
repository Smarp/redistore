
package redistore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRedistore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redistore Suite")
}
