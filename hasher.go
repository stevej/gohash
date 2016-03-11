package gohash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

// ...
const (
	AllowedOnion = "abcdefghijklmnopqrstuvwxyz234567"
)

var (
	algos = map[string]int{
		"md5":        128,
		"sha1":       160,
		"sha224":     224,
		"sha256":     256,
		"sha384":     384,
		"sha512":     512,
		"sha512-224": 224,
		"sha512-256": 256,
		"sha3-224":   224,
		"sha3-256":   256,
		"sha3-384":   384,
		"sha3-512":   512,
	}
)

// Hasher ...
type Hasher struct {
	algo        string
	suffix      string
	expected    []byte
	minLength   int
	maxLength   int
	allowedKeys []byte
	reverse     bool

	// for runtime stats
	buffer []byte
	try    uint64
	tick   uint64
}

// NewHasher returns a new Hasher
func NewHasher() *Hasher {
	return &Hasher{}
}

// Algo sets the hash algorithm ("sha1", "sha512")
func (h *Hasher) Algo(algo string) {
	algo = strings.Replace(algo, "_", "-", -1)
	h.algo = algo
}

// ExpectedHash sets the expected hash
func (h *Hasher) ExpectedHash(expected string) {
	tmp := hexStringToBytes(expected)
	h.expected = tmp[:]
}

// Length sets the length of key to find
func (h *Hasher) Length(len int) {
	h.minLength = len
	h.maxLength = len
}

// Reverse sets wether to do sequential find in reverse order
func (h *Hasher) Reverse(b bool) {
	h.reverse = b
}

// Suffix sets a fixed suffix
func (h *Hasher) Suffix(s string) { h.suffix = s }

// MinLength sets min length of key to find
func (h *Hasher) MinLength(len int) { h.minLength = len }

// MaxLength sets max length of key to find
func (h *Hasher) MaxLength(len int) { h.maxLength = len }

// AllowedKeys sets the allowed keys
func (h *Hasher) AllowedKeys(s string) {
	h.allowedKeys = strToDistinctByteSlice(s)
}

// GetAllowedKeys ...
func (h *Hasher) GetAllowedKeys() string { return string(h.allowedKeys) }

func (h *Hasher) verify() error {

	if len(h.allowedKeys) == 0 {
		return fmt.Errorf("allowedKeys unset")
	}

	if h.minLength == 0 {
		return fmt.Errorf("minLength unset")
	}

	if len(h.algo) == 0 {
		return fmt.Errorf("algo unset")
	}

	keyBitSize := len(h.expected) * 8
	expectedBitSize := len(h.expected) * 8

	if requiredBitSize, ok := algos[h.algo]; ok {
		if keyBitSize != requiredBitSize {
			return fmt.Errorf("expectedHash is wrong size, should be %d bit, is %d",
				requiredBitSize, expectedBitSize)
		}
	} else {
		return fmt.Errorf("unknown algo %s", h.algo)
	}

	return nil
}

func (h *Hasher) equals() bool {

	if h.algo == "md5" && byte16ArrayEquals(md5.Sum(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha1" && byte20ArrayEquals(sha1.Sum(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha224" && byte28ArrayEquals(sha256.Sum224(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha256" && byte32ArrayEquals(sha256.Sum256(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha384" && byte48ArrayEquals(sha512.Sum384(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha512" && byte64ArrayEquals(sha512.Sum512(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha512-224" && byte28ArrayEquals(sha512.Sum512_224(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha512-256" && byte32ArrayEquals(sha512.Sum512_256(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha3-224" && byte28ArrayEquals(sha3.Sum224(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha3-256" && byte32ArrayEquals(sha3.Sum256(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha3-384" && byte48ArrayEquals(sha3.Sum384(h.buffer), h.expected) {
		return true
	}

	if h.algo == "sha3-512" && byte64ArrayEquals(sha3.Sum512(h.buffer), h.expected) {
		return true
	}

	return false
}

// FindSequential calcs all possible combinations of keys of given length
func (h *Hasher) FindSequential() (string, error) {

	if err := h.verify(); err != nil {
		return "", err
	}

	h.buffer = make([]byte, h.minLength)

	firstAllowedKey := h.allowedKeys[0]
	lastAllowedKey := h.allowedKeys[len(h.allowedKeys)-1]

	// create initial mutation
	for x := 0; x < h.minLength; x++ {
		if h.reverse {
			h.buffer[x] = lastAllowedKey
		} else {
			h.buffer[x] = firstAllowedKey
		}
	}

	h.buffer = append(h.buffer, h.suffix...)

	go h.statusReport()

	for {

		if h.equals() {
			return string(h.buffer), nil
		}

		// update mutation
		for roller := h.minLength - 1; roller >= 0; roller-- {
			if h.reverse {
				if h.buffer[roller] == firstAllowedKey {
					h.buffer[roller] = lastAllowedKey
					continue
				} else {
					h.buffer[roller] = h.prevValueFor(h.buffer[roller])
					break
				}
			} else {
				if h.buffer[roller] == lastAllowedKey {
					h.buffer[roller] = firstAllowedKey
					continue
				} else {
					h.buffer[roller] = h.nextValueFor(h.buffer[roller])
					break
				}
			}
		}
		h.try++
	}
}

// FindRandom uses random brute force to attempt to find by luck
func (h *Hasher) FindRandom() (string, error) {

	if h.reverse {
		return "", fmt.Errorf("reverse and random dont mix")
	}

	if err := h.verify(); err != nil {
		return "", err
	}

	h.buffer = make([]byte, h.minLength)

	firstAllowedKey := h.allowedKeys[0]
	allowedKeysLen := len(h.allowedKeys)

	// create initial mutation
	for x := 0; x < h.minLength; x++ {
		h.buffer[x] = firstAllowedKey
	}

	h.buffer = append(h.buffer, h.suffix...)

	go h.statusReport()

	for {
		if h.equals() {
			return string(h.buffer), nil
		}

		// update mutation of first letters
		for roller := 0; roller < h.minLength; roller++ {
			h.buffer[roller] = h.allowedKeys[rand.Intn(allowedKeysLen)]
		}
		h.try++
	}
}

func (h *Hasher) statusReport() {

	for {
		time.Sleep(1 * time.Second)
		h.tick++
		avg := h.try / h.tick

		fmt.Printf("%s ~%d/s %s\n", h.algo, avg, string(h.buffer))
	}
}

func (h *Hasher) nextValueFor(b byte) byte {

	next := false
	for _, x := range h.allowedKeys {
		if next == true {
			return x
		}
		if x == b {
			next = true
		}
	}
	return '0'
}

func (h *Hasher) prevValueFor(b byte) byte {

	prev := h.allowedKeys[0]
	for _, x := range h.allowedKeys {
		if x == b {
			return prev
		}
		prev = x
	}
	return prev
}
