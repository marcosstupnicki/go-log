package golog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Level
	}{
		{name: "debug", in: "debug", want: DebugLevel},
		{name: "info", in: "info", want: InfoLevel},
		{name: "warn", in: "warn", want: WarnLevel},
		{name: "warning alias", in: "warning", want: WarnLevel},
		{name: "error", in: "error", want: ErrorLevel},
		{name: "trim and case", in: " WARN ", want: WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.in)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseLevelRejectsUnknownLevel(t *testing.T) {
	_, err := ParseLevel("verbose")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "debug, info, warn, error")
}

func TestMustParseLevelPanicsOnUnknownLevel(t *testing.T) {
	assert.Panics(t, func() { _ = MustParseLevel("verbose") })
}
