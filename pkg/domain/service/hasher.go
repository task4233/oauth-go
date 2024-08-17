package service

import (
	"bytes"
	"context"
	"crypto/sha256"

	"github.com/task4233/oauth/pkg/repository"
)

var _ repository.Hasher = (*Sha256Hasher)(nil)

type Sha256Hasher struct {
	fixedKey string
}

func NewSha256Hasher(fixedKey string) *Sha256Hasher {
	return &Sha256Hasher{
		fixedKey: fixedKey,
	}
}

func (s *Sha256Hasher) Compare(_ context.Context, hash, data []byte) (bool, error) {
	target := sha256.Sum256([]byte(string(data) + s.fixedKey))
	return bytes.Equal(hash, target[:]), nil
}

func (s *Sha256Hasher) Hash(_ context.Context, data []byte) ([]byte, error) {
	target := sha256.Sum256([]byte(string(data) + s.fixedKey))
	return target[:], nil
}
