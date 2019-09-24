package gentoo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGentoo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gentoo Suite")
}
