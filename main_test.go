package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/BytecodeAgency/import-boundary-checker/logging"
	"github.com/BytecodeAgency/import-boundary-checker/runner"
	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	rootDir, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		dir       string
		shouldErr bool
	}{
		{"go-invalid-1", true},
		{"go-invalid-2", true},
		{"go-invalid-3", true},
		{"go-invalid-4", true},
		{"go-invalid-5", true},
		{"go-valid-1", false},
		{"go-valid-2", false},
		{"go-valid-3", false},
	}
	for _, test := range tests {

		// cd into test dir
		err := os.Chdir("./examples/" + test.dir)
		assert.NoError(t, err)

		// Load config file
		abs, err := filepath.Abs(".importrules")
		assert.NoError(t, err)
		configFile, err := ioutil.ReadFile(abs)
		assert.NoError(t, err)
		config := string(configFile)

		// Create and run runner
		// TODO: Add automated end-to-end tests
		logger := logging.Logger{Verbose: false}
		r := runner.New(config, &logger)
		got := r.Run()

		// Check if we got what we expected
		assert.Equalf(t, test.shouldErr, got, "example directory %s", test.dir)

		// Change back to parent directory
		err = os.Chdir(rootDir)
		assert.NoError(t, err)
	}
}
