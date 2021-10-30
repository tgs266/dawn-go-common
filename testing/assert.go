package testing

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func StatusCodeEqual(t *testing.T, expected int, actual int) {
	assert.Equal(t, expected, actual, "status code must be "+strconv.Itoa(actual))
}
