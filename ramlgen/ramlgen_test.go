package main

import (
	"testing"

	"github.com/EconomistDigitalSolutions/ramlapi"
)

func TestValidRAML(t *testing.T) {
	_, err := ramlapi.ProcessRAML("../fixtures/valid.raml")
	if err != nil {
		t.Fatalf("Expected ramlgen to parse valid RAML, got %v\n", err)
	}
}

func TestInvalidRAML(t *testing.T) {
	_, err := ramlapi.ProcessRAML("../fixtures/invalid.raml")
	if err == nil {
		t.Fatal("Expected ramlgen to abort on in valid RAML\n")
	}
}

func TestValidRAMLData(t *testing.T) {
	api, _ := ramlapi.ProcessRAML("../fixtures/valid.raml")
	if api.Title != "gref" {
		t.Fatalf("Expected API title to be parsed as gref, got %s\n", api.Title)
	}
}
