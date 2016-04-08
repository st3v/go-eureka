package jolt_test

import (
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJolt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jolt")
}

func removeIdendation(data []byte) []byte {
	r := regexp.MustCompile("\\n\\s*")
	return r.ReplaceAll(data, []byte{})
}
