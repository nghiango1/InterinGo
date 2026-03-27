package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func isIigFile(fullName string) (string, bool) {
	splited := strings.Split(fullName, ".")
	if len(splited) < 2 {
		return "", false
	}

	extension := splited[len(splited)-1]
	if extension != "iig" {
		return "", false
	}

	fileName := strings.Join(splited[0:len(splited)-1], "")
	return fileName, true
}

// Hack so that testing flag can be init first
var _ = func() bool {
	testing.Init()
	return true
}()

func findProjectRoot() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }

    // Traverse up the directory tree to find go.mod
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir, nil // Return the directory containing go.mod
        }
        parentDir := filepath.Dir(dir)
        if parentDir == dir { // Reached the root directory, stop
            break
        }
        dir = parentDir
    }
    return "", os.ErrNotExist // go.mod not found
}
