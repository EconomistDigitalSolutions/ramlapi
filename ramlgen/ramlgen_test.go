package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/EconomistDigitalSolutions/ramlapi"
)

var output = "/%s/test_gen_%d"

func TestGenerate(t *testing.T) {
	api, _ := ramlapi.Process("../fixtures/valid.raml")
	currentOutput := fmt.Sprintf(output, os.TempDir(), int32(time.Now().Unix()))
	generate(api, currentOutput)
	_, err := os.Open(currentOutput)
	if err != nil {
		t.Fatalf("Expected output file to exist, got %v\n", err)
	}
	os.Remove(currentOutput)
}
