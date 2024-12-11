package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

type Crypto interface {
	Generate([]byte) []byte
	GenerateBase64String([]byte) string
	CompareValueToHash(string, []byte) error
}
type crypto struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// salt
	salt []byte
}

func NewCrypto(time uint32, salt []byte, memory uint32, threads uint8, keyLen uint32) Crypto {
	return &crypto{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		salt:    salt,
	}
}

// https://stackoverflow.com/questions/39481826/generate-6-digit-verification-code-with-golang
func EncodeToString(max int) string {
	var randomCodeTable [10]byte = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = randomCodeTable[int(b[i])%len(randomCodeTable)]
	}
	return string(b)
}

func (h *crypto) Generate(value []byte) []byte {
	return argon2.IDKey(value, h.salt, h.time, h.memory, h.threads, h.keyLen)
}

func (h *crypto) GenerateBase64String(value []byte) string {
	return base64.StdEncoding.EncodeToString(h.Generate(value))
}

func (h *crypto) CompareValueToHash(value string, hash []byte) error {
	valueHash := h.Generate([]byte(value))
	if !bytes.Equal(valueHash, hash) {
		return errors.New("does not match")
	}
	return nil
}
