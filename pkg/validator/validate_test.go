package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInput struct {
	Password string `json:"password" validate:"required,password"`
	Email    string `json:"email" validate:"required,email"`
}

func TestCustomValidator_Validate_Valid(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "Aa1!aaaa",
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	assert.NoError(t, err)
}

func TestCustomValidator_Validate_PasswordTooShort(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "Aa1!a", // too short
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "must be between"))
	}
}

func TestCustomValidator_Validate_PasswordNoLower(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "AAAAA1!A",
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "lowercase"))
	}
}

func TestCustomValidator_Validate_PasswordNoUpper(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "aaaaa1!a",
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "uppercase"))
	}
}

func TestCustomValidator_Validate_PasswordNoDigit(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "Aaaaaa!A",
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "digit"))
	}
}

func TestCustomValidator_Validate_PasswordNoSymbol(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "Aaaaaa1A",
		Email:    "a@b.co",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "special character"))
	}
}

func TestCustomValidator_Validate_InvalidEmail(t *testing.T) {
	cv := NewCustomValidator()
	in := testInput{
		Password: "Aa1!aaaa",
		Email:    "invalid",
	}
	err := cv.Validate(in)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "must be a valid email address"))
	}
}
