package config

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestFromFile(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "examplesite")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	defer os.RemoveAll(temporaryDirectory)

	filePath := path.Join(temporaryDirectory, "config.yaml")
	contents := `---
title: Test
description: This is a test.
database_file: test/example.sqlite
items_per_page: 100
`
	configFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_, err = configFile.WriteString(contents)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	configFile.Close()

	expected := Config{
		"Test",
		"This is a test.",
		"test/example.sqlite",
		100,
	}

	result, err := FromFile(filePath)
	if !reflect.DeepEqual(result, &expected) {
		t.Fatalf("Received:\n%+v\nExpected:\n%+v", result, expected)
	}
}
