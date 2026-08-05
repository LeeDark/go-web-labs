package main

import (
	"testing"

	"go.uber.org/fx"
)

func TestApplicationOptionsValidate(t *testing.T) {
	if err := fx.ValidateApp(applicationOptions(":0")...); err != nil {
		t.Fatalf("Fx application validation failed: %v", err)
	}
}
