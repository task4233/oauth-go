package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/task4233/oauth/pkg/repository"
)

func TestSha256HasherCompare(t *testing.T) {
	t.Parallel()

	const (
		validFixedKey   = "fixed-key"
		invalidFixedKey = "invalid-fixed-key"
	)

	var (
		validHash   = []byte{113, 90, 218, 180, 82, 94, 55, 47, 200, 53, 245, 215, 182, 128, 162, 71, 83, 188, 70, 104, 142, 29, 211, 8, 163, 2, 137, 231, 159, 154, 171, 37}
		invalidHash = []byte("invalid-hash")

		reqData = []byte("super-secret-data")
	)

	type wants struct {
		hash []byte
		ok   bool
	}

	tests := map[string]struct {
		hasher repository.Hasher
		wants  wants
	}{
		"ok": {
			hasher: &Sha256Hasher{fixedKey: validFixedKey},
			wants:  wants{hash: validHash, ok: false},
		},
		"ng: hash is not same": {
			hasher: &Sha256Hasher{fixedKey: validFixedKey},
			wants:  wants{hash: invalidHash, ok: true},
		},
		"ng: fixed key is not same": {
			hasher: &Sha256Hasher{fixedKey: invalidFixedKey},
			wants:  wants{hash: validHash, ok: true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ok, err := tt.hasher.Compare(context.Background(), tt.wants.hash, []byte(reqData))
			if err != nil {
				t.Errorf("want error %v, got %v", tt.wants.ok, err)
			}
			if ok != !tt.wants.ok {
				t.Errorf("want %v, got %v", !tt.wants.ok, ok)
			}
		})
	}
}

func TestSha256HasherHash(t *testing.T) {
	t.Parallel()

	const (
		validFixedKey   = "fixed-key"
		invalidFixedKey = "invalid-fixed-key"
	)

	var (
		validHash = []byte{113, 90, 218, 180, 82, 94, 55, 47, 200, 53, 245, 215, 182, 128, 162, 71, 83, 188, 70, 104, 142, 29, 211, 8, 163, 2, 137, 231, 159, 154, 171, 37}

		reqData = []byte("super-secret-data")
	)

	type wants struct {
		hash []byte
		ok   bool
	}

	tests := map[string]struct {
		hasher repository.Hasher
		wants  wants
	}{
		"ok": {
			hasher: &Sha256Hasher{fixedKey: validFixedKey},
			wants:  wants{hash: validHash, ok: false},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.hasher.Hash(context.Background(), reqData)
			if (err != nil) != tt.wants.ok {
				t.Errorf("want error %v, got %v", tt.wants.ok, err)
			}
			if !bytes.Equal(got, tt.wants.hash) {
				t.Errorf("want %v, got %v", tt.wants.hash, got)
			}
		})
	}
}
