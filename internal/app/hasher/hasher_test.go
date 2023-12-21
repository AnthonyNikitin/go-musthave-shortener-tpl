package hasher

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetShortLink(t *testing.T) {

	tests := []struct {
		name         string
		input        string
		wantResponse string
		wantError    error
	}{
		{
			name:         "positive test",
			input:        "https://practicum.yandex.ru/",
			wantResponse: "CKj87ajs",
			wantError:    nil,
		},
		{
			name:         "negative test",
			input:        "",
			wantResponse: "",
			wantError:    errors.New("input should not be empty"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shortLink, err := GetShortLink(test.input)

			assert.Equal(t, test.wantError, err)
			assert.Equal(t, test.wantResponse, shortLink)
		})
	}
}
