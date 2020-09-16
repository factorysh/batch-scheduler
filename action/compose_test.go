package action

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCompose(t *testing.T) {
	tests := map[string]struct {
		input      string
		shouldFail bool
	}{
		"valid":   {input: "../compose-samples/valid-echo.yml", shouldFail: false},
		"invalid": {input: "../compose-samples/invalid-echo.yml", shouldFail: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			input, err := ioutil.ReadFile(tc.input)
			assert.NoError(t, err)

			c := NewCompose(input)

			err = c.Parse()
			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func NewValidCompose() (*Compose, error) {

	input, err := ioutil.ReadFile("../compose-samples/valid-echo.yml")
	if err != nil {
		return nil, err
	}

	c := NewCompose(input)

	err = c.Parse()
	if err != nil {
		return nil, err
	}

	return &c, nil

}

func TestRecompose(t *testing.T) {

	expected := `services:
  hello:
    command: echo world
    image: busybox:latest
version: "3"
`

	c, err := NewValidCompose()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	output, err := c.Recompose()
	assert.NoError(t, err)

	assert.Equal(t, expected, output)

}

func TestEnsureBin(t *testing.T) {
	tests := map[string]struct {
		input      string
		shouldFail bool
	}{
		"valid":   {input: "docker-compose", shouldFail: false},
		"invalid": {input: "dckr-cmps", shouldFail: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ensureBin(tc.input)
			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestRun(t *testing.T) {
	c, err := NewValidCompose()
	assert.NoError(t, err)

	err = c.Run("uuid")
	assert.NoError(t, err)

}
