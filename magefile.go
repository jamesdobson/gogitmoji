// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

/*
 * Targets
 */

// Lints the code.
func Lint() {
	sh.Exec(
		nil, os.Stdout, os.Stderr,
		"golangci-lint",
		"run",
		"--enable=bodyclose",
		"--enable=dupl",
		"--enable=funlen",
		"--enable=gocognit",
		"--enable=goconst",
		"--enable=gocritic",
		"--enable=gocyclo",
		"--enable=gofmt",
		"--enable=goimports",
		"--enable=golint",
		"--enable=gosec",
		"--enable=interfacer",
		"--enable=maligned",
		"--enable=prealloc",
		"--enable=scopelint",
		"--enable=unparam",
		"--enable=unconvert",
		"--enable=whitespace",
	)
}