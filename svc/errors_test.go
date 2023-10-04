package svc

import (
	"errors"
	"testing"
)

func TestNewErrValidation(t *testing.T) {
	err := NewErrValidation("invalid input")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "invalid input" {
		t.Errorf("expected error message 'invalid input', got '%s'", err.Error())
	}
}

func TestNewErrServerError(t *testing.T) {
	cause := errors.New("database error")
	err := NewErrServerError("server error", cause)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "server error" {
		t.Errorf("expected error message 'server error', got '%s'", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Errorf("expected error cause '%v', got '%v'", cause, err)
	}
}

func TestNewErrConflict(t *testing.T) {
	err := NewErrConflict("resource conflict")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "resource conflict" {
		t.Errorf("expected error message 'resource conflict', got '%s'", err.Error())
	}
}

func TestNewErrNotFound(t *testing.T) {
	err := NewErrNotFound("resource not found")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "resource not found" {
		t.Errorf("expected error message 'resource not found', got '%s'", err.Error())
	}
}
