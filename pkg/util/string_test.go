package util_test

import (
	"testing"

	"github.com/gmcoringa/tswitch/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestIsBlankWithEmptyStringShouldReturnTrue(testing *testing.T) {
	result := util.IsBlank("")

	assert.True(testing, result)
}

func TestIsBlankWithBlankStringShouldReturnTrue(testing *testing.T) {
	result := util.IsBlank("  ")

	assert.True(testing, result)
}

func TestIsBlankWithNonBlankStringShouldReturnFalse(testing *testing.T) {
	result := util.IsBlank(" a ")

	assert.False(testing, result)
}

func TestIsNotBlankWithEmptyStringShouldReturnFalse(testing *testing.T) {
	result := util.IsNotBlank("")

	assert.False(testing, result)
}

func TestIsNotBlankWithNonBlankStringShouldReturnTrue(testing *testing.T) {
	result := util.IsNotBlank(" a ")

	assert.True(testing, result)
}
