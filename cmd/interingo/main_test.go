package main

import (
	"bytes"
	"interingo/pkg/repl"
	"log"
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

// This evaluating all `test/*.iig` file and compare it with the desier
// `test/result/*.out` file content
func TestMain(t *testing.T) {
	// Get the project root directory (where go.mod is located)
	rootDir, err := findProjectRoot()
	if err != nil {
		log.Fatal(err)
	}

	// Construct asset path assuming the test directory is at the project root
	testAssetsDir := filepath.Join(rootDir, "test")
	f, err := os.Open(testAssetsDir)
	if err != nil {
		t.Errorf("Test directory read error, is it available, error code: %v\n", err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		t.Errorf("Test directory read error, is it available, error code: %v\n", err)
	}

	for _, v := range files {
		buf := new(bytes.Buffer)
		if v.IsDir() {
			continue
		}

		fileName, ok := isIigFile(v.Name())

		if !ok {
			continue
		}

		inputFileDir := filepath.Join(testAssetsDir, v.Name())
		inputFileContent, err := os.ReadFile(inputFileDir)
		if err != nil {
			t.Errorf("File read error, error code: %v\n", err)
		}

		repl := repl.NewRepl(nil, os.Stdin, buf)
		repl.Handle(string(inputFileContent))

		outputFileName := fileName + ".out"
		outputFileNameDir := filepath.Join(testAssetsDir, "/result/", outputFileName)
		outputFileContent, err := os.ReadFile(outputFileNameDir)
		if err != nil {
			t.Errorf("File read error, recheck output file %s location\n", outputFileName)
		}

		for i, outByte := range outputFileContent {
			b, err := buf.ReadByte()
			if err != nil {
				break
			}
			if outByte != b {
				t.Errorf("Result of %s not match output file, expect \"%c\" at %d'th output but got \"%c\" instead", v.Name(), outByte, i, b)
			}
		}
	}
}
