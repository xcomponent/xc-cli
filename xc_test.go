package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestXcCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "XC CLI Tests")
}
