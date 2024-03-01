package main

import (
	"bytes"
	"interingo/repl"
	"os"
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

// This evaluating all `test/*.iig` file and compare it with the desier
// `test/result/*.out` file content
func TestMain(t *testing.T) {
	f, err := os.Open("test")
	if err != nil {
		t.Errorf("Test directory read error, is it available, error code: %v\n", err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		t.Errorf("Test directory read error, is it available, error code: %v\n", err)
	}

	buf := new(bytes.Buffer)
	for _, v := range files {
		if v.IsDir() {
			continue
		}

		fileName, ok := isIigFile(v.Name())

		if !ok {
			continue
		}

		inputFileContent, err := os.ReadFile("test/" + v.Name())
		if err != nil {
			t.Errorf("File read error, error code: %v\n", err)
		}

		repl.Handle(string(inputFileContent), buf)

		outputFileName := fileName + ".out"
		outputFileContent, err := os.ReadFile("test/result/" + outputFileName)
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
