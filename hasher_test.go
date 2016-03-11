package gohash

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexStringToBytes(t *testing.T) {

	assert.Equal(t, []byte("hej"), hexStringToBytes("68656a"))
}

func TestHashSha1(t *testing.T) {

	hasher := NewHasher()

	hasher.Algo("sha1")
	hasher.AllowedKeys("tom")
	hasher.ExpectedHash("9af7d87edaba03e23f6dbdaed29101ee1291c8a6")
	hasher.Length(5)

	res, err := hasher.FindSequential()
	assert.Equal(t, nil, err)
	assert.Equal(t, "motom", string(res))
}

func TestHashSha512(t *testing.T) {

	hasher := NewHasher()

	hasher.Algo("sha512")
	hasher.AllowedKeys("mota")
	hasher.ExpectedHash("2b4df6c7b86d49a71c5ad6c1ffc2e85fde618c69400d2a0ccabb8dd12df4ae2584103b6947379742c0bc11a4e81ad4a3a832c11a734bf8ae5f8a8af9317a4c03")
	hasher.Length(4)

	assert.Equal(t, "amot", hasher.GetAllowedKeys()) // make sure they are sorted alphabetically

	res, err := hasher.FindSequential()
	assert.Equal(t, nil, err)
	assert.Equal(t, "atom", string(res))
}

func TestHashSha512OnionAddress(t *testing.T) {

	hasher := NewHasher()

	hasher.Algo("sha512")
	hasher.AllowedKeys("abcdefghijklmnopqrstuvwxyz234567")
	hasher.Suffix(".onion")
	hasher.ExpectedHash("f07be23625ad049e9c44d9d2a8088d3902f5dbbd3f16a1469c34051d5987c5859fc1eeb0127764ad1ba1de4da51297002baaa1b41f3e259d54b135434d8851cc")
	hasher.Length(16)

	res, err := hasher.FindSequential()
	assert.Equal(t, nil, err)
	assert.Equal(t, "222222222222222f.onion", string(res))
}

func TestHashSha512OnionAddressRandomFind(t *testing.T) {

	rand.Seed(123)

	hasher := NewHasher()

	hasher.Algo("sha512")
	hasher.AllowedKeys("abcdefghijklmnopqrstuvwxyz234567")
	hasher.Suffix(".onion")
	hasher.ExpectedHash("bbc3581fa536cf90d95b60d226495d38257d73e971b3193cc3fd532338caba7710966c5c91eddc8c1193e9cf401db94cb7c16205f064b6c45e3320d8c5d0b5f3")
	hasher.Length(16)

	res, err := hasher.FindRandom()
	assert.Equal(t, nil, err)
	assert.Equal(t, "2gl57brnwcjqmaua.onion", string(res))
}

// benchmarks given key length and print a prediction based on it
func BenchmarkSha1Speed(*testing.B) {

	hasher := NewHasher()

	hasher.Algo("sha1")
	hasher.AllowedKeys("580%(=QWI+qwi*Nn")
	hasher.ExpectedHash("0000000000000000000000000000000000000000")
	hasher.Length(5)

	hasher.FindSequential()
}
