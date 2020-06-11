// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var goEnv = map[string]string{ // nolint: gochecknoglobals,unused
	"GOOS":        "linux",
	"GOARCH":      "amd64",
	"CGO_ENABLED": "0",
}

// Clean removes any build artifacts.
func Clean() error { // nolint: deadcode,unused
	if err := os.RemoveAll("dist"); err != nil {
		return mg.Fatalf(1, "error removing dist directory: %w", err)
	}

	return nil
}

// Build will compile the project.
func Build() error { // nolint: deadcode,unused
	if err := sh.RunWith(goEnv,
		"go", "build", "-o", "dist/asmith", "./cmd/asmith",
	); err != nil {
		return mg.Fatalf(1, "error building asmith: %w", err)
	}

	return nil
}

// Lint runs static code analysis.
func Lint() error { // nolint: deadcode,unused
	if err := sh.RunWith(goEnv, "golangci-lint", "run"); err != nil {
		return mg.Fatalf(1, "error building asmith: %w", err)
	}

	return nil
}
