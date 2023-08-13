package api

import "testing"

func TestAuthorize(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
	}{
		"success": {},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_ = tt
		})
	}
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
	}{
		"success": {},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_ = tt
		})
	}
}

func TestToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{}{
		"success": {},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_ = tt
		})
	}
}
