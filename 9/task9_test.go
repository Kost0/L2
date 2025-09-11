package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack_Success(t *testing.T) {
	str := "a4b5"
	res, err := unpack(str)

	expect := "aaaabbbbb"

	assert.Equal(t, expect, res)
	assert.NoError(t, err)
}

func TestUnpack_SuccessWithEscape(t *testing.T) {
	str := "a4b/53"
	res, err := unpack(str)

	expect := "aaaab555"

	assert.Equal(t, expect, res)
	assert.NoError(t, err)
}

func TestUnpack_SuccessWithoutDigitAtTheEnd(t *testing.T) {
	str := "a4b"
	res, err := unpack(str)

	expect := "aaaab"

	assert.Equal(t, res, expect)
	assert.NoError(t, err)
}

func TestUnpack_SuccessWithoutDigitAtTheEnd2(t *testing.T) {
	str := "a4b/5"
	res, err := unpack(str)

	expect := "aaaab5"

	assert.Equal(t, res, expect)
	assert.NoError(t, err)
}

func TestUnpack_SuccessDifferentNumbersClose(t *testing.T) {
	str := "a/45"
	res, err := unpack(str)

	expect := "a44444"

	assert.Equal(t, res, expect)
	assert.NoError(t, err)
}

func TestUnpack_SuccessWithoutDigits(t *testing.T) {
	str := "ab"
	res, err := unpack(str)

	expect := "ab"

	assert.Equal(t, res, expect)
	assert.NoError(t, err)
}

func TestUnpack_SuccessEmptyString(t *testing.T) {
	str := ""
	res, err := unpack(str)

	expect := ""

	assert.Equal(t, res, expect)
	assert.NoError(t, err)
}

func TestUnpack_FailOnlyDigits(t *testing.T) {
	str := "45"
	res, err := unpack(str)

	assert.Equal(t, res, "")
	assert.Error(t, err)
}
