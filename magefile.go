//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

/*
 * Targets
 */

// Lint checks the code for common errors.
func Lint() error {
	_, err := sh.Exec(
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
		"--enable=gosec",
		"--enable=prealloc",
		"--enable=revive",
		"--enable=unparam",
		"--enable=unconvert",
		"--enable=whitespace",
	)

	return err
}
