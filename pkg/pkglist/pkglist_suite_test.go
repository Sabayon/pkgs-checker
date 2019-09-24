package pkglist_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommons(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkglist Suite")
}
