package gohash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"sort"

	"github.com/cxmcc/tiger"
	"github.com/dchest/blake256"
	"github.com/dchest/blake2b"
	"github.com/dchest/blake2s"
	"github.com/dchest/blake512"
	"github.com/dchest/siphash"
	"github.com/dchest/skein"
	"github.com/htruong/go-md2"
	"github.com/jzelinskie/whirlpool"
	"github.com/stargrave/gogost/gost341194"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Calculator is used to calculate hash of input cleartext
type Calculator struct {
	data []byte
}

// NewCalculator creates a new Calculator
func NewCalculator(data []byte) *Calculator {

	return &Calculator{
		data: data,
	}
}

var (
	algos = map[string]int{
		"adler32":      32,
		"blake224":     224,
		"blake256":     256,
		"blake384":     384,
		"blake512":     512,
		"blake2b-512":  512,
		"blake2s-256":  256,
		"crc32":        32,
		"crc32c":       32,
		"crc32k":       32,
		"gost":         256,
		"md2":          128,
		"md4":          128,
		"md5":          128,
		"ripemd160":    160,
		"sha1":         160,
		"sha224":       224,
		"sha256":       256,
		"sha384":       384,
		"sha512":       512,
		"sha512-224":   224,
		"sha512-256":   256,
		"sha3-224":     224,
		"sha3-256":     256,
		"sha3-384":     384,
		"sha3-512":     512,
		"shake128-256": 256,
		"shake256-512": 512,
		"siphash-2-4":  64,
		"skein512-256": 256,
		"skein512-512": 512,
		"tiger192":     192,
		"whirlpool":    512,
	}

	checksummers = map[string]func(*[]byte) *[]byte{
		"adler32":      adler32Sum,
		"blake224":     blake224Sum,
		"blake256":     blake256Sum,
		"blake384":     blake384Sum,
		"blake512":     blake512Sum,
		"blake2b-512":  blake2b_512Sum,
		"blake2s-256":  blake2s_256Sum,
		"crc32":        crc32Sum,
		"crc32c":       crc32cSum,
		"crc32k":       crc32kSum,
		"fnv1-32":      fnv1_32Sum,
		"fnv1a-32":     fnv1a_32Sum,
		"fnv1-64":      fnv1_64Sum,
		"fnv1a-64":     fnv1a_64Sum,
		"gost":         gostSum,
		"md2":          md2Sum,
		"md4":          md4Sum,
		"md5":          md5Sum,
		"ripemd160":    ripemd160Sum,
		"sha1":         sha1Sum,
		"sha224":       sha224Sum,
		"sha256":       sha256Sum,
		"sha384":       sha384Sum,
		"sha512":       sha512Sum,
		"sha512-224":   sha512_224Sum,
		"sha512-256":   sha512_256Sum,
		"sha3-224":     sha3_224Sum,
		"sha3-256":     sha3_256Sum,
		"sha3-384":     sha3_384Sum,
		"sha3-512":     sha3_512Sum,
		"shake128-256": shake128_256Sum,
		"shake256-512": shake256_512Sum,
		"siphash-2-4":  siphash2_4Sum,
		"skein512-256": skein512_256Sum,
		"skein512-512": skein512_512Sum,
		"tiger192":     tiger192Sum,
		"whirlpool":    whirlpoolSum,
	}
)

// Sum returns the checksum
func (c *Calculator) Sum(algo string) *[]byte {

	algo = resolveAlgoAliases(algo)

	if checksum, ok := checksummers[algo]; ok {
		return checksum(&c.data)
	}
	return nil
}

// AvailableHashes returns the available hash id's
func AvailableHashes() []string {

	res := []string{}

	for key := range checksummers {
		res = append(res, key)
	}

	sort.Strings(res)
	return res
}

func resolveAlgoAliases(s string) string {

	// "tiger" is used by rhash, sphsum
	if s == "tiger" {
		return "tiger192"
	}

	// "skein256" is used in sphsum
	if s == "skein256" {
		return "skein512-256"
	}

	// "skein512" is used in sphsum
	if s == "skein512" {
		return "skein512-256"
	}

	return s
}

func adler32Sum(b *[]byte) *[]byte {
	i := adler32.Checksum(*b)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return &bs
}

func blake224Sum(b *[]byte) *[]byte {
	w := blake256.New224()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func blake256Sum(b *[]byte) *[]byte {
	w := blake256.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func blake384Sum(b *[]byte) *[]byte {
	w := blake512.New384()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func blake512Sum(b *[]byte) *[]byte {
	w := blake512.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func blake2b_512Sum(b *[]byte) *[]byte {
	x := blake2b.Sum512(*b)
	res := x[:]
	return &res
}

func blake2s_256Sum(b *[]byte) *[]byte {
	x := blake2s.Sum256(*b)
	res := x[:]
	return &res
}

func crc32Sum(b *[]byte) *[]byte {
	i := crc32.ChecksumIEEE(*b)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return &bs
}

func crc32cSum(b *[]byte) *[]byte {
	tbl := crc32.MakeTable(crc32.Castagnoli)
	i := crc32.Checksum(*b, tbl)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return &bs
}

func crc32kSum(b *[]byte) *[]byte {
	tbl := crc32.MakeTable(crc32.Koopman)
	i := crc32.Checksum(*b, tbl)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return &bs
}

func fnv1_32Sum(b *[]byte) *[]byte {
	w := fnv.New32()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func fnv1a_32Sum(b *[]byte) *[]byte {
	w := fnv.New32a()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func fnv1_64Sum(b *[]byte) *[]byte {
	w := fnv.New64()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func fnv1a_64Sum(b *[]byte) *[]byte {
	w := fnv.New64a()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func gostSum(b *[]byte) *[]byte {
	h := gost341194.New(gost341194.SboxDefault)
	h.Write(*b)
	res := h.Sum(nil)
	return &res
}

func md2Sum(b *[]byte) *[]byte {
	w := md2.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func md4Sum(b *[]byte) *[]byte {
	w := md4.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func md5Sum(b *[]byte) *[]byte {
	x := md5.Sum(*b)
	res := x[:]
	return &res
}

func ripemd160Sum(b *[]byte) *[]byte {
	w := ripemd160.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func sha1Sum(b *[]byte) *[]byte {
	x := sha1.Sum(*b)
	res := x[:]
	return &res
}

func sha224Sum(b *[]byte) *[]byte {
	x := sha256.Sum224(*b)
	res := x[:]
	return &res
}

func sha256Sum(b *[]byte) *[]byte {
	x := sha256.Sum256(*b)
	res := x[:]
	return &res
}

func sha384Sum(b *[]byte) *[]byte {
	x := sha512.Sum384(*b)
	res := x[:]
	return &res
}

func sha512Sum(b *[]byte) *[]byte {
	x := sha512.Sum512(*b)
	res := x[:]
	return &res
}

func sha512_224Sum(b *[]byte) *[]byte {
	x := sha512.Sum512_224(*b)
	res := x[:]
	return &res
}

func sha512_256Sum(b *[]byte) *[]byte {
	x := sha512.Sum512_256(*b)
	res := x[:]
	return &res
}

func sha3_224Sum(b *[]byte) *[]byte {
	x := sha3.Sum224(*b)
	res := x[:]
	return &res
}

func sha3_256Sum(b *[]byte) *[]byte {
	x := sha3.Sum256(*b)
	res := x[:]
	return &res
}

func sha3_384Sum(b *[]byte) *[]byte {
	x := sha3.Sum384(*b)
	res := x[:]
	return &res
}

func sha3_512Sum(b *[]byte) *[]byte {
	x := sha3.Sum512(*b)
	res := x[:]
	return &res
}

func shake128_256Sum(b *[]byte) *[]byte {
	res := make([]byte, 32)
	sha3.ShakeSum128(res, *b)
	return &res
}

func shake256_512Sum(b *[]byte) *[]byte {
	res := make([]byte, 64)
	sha3.ShakeSum256(res, *b)
	return &res
}

func siphash2_4Sum(b *[]byte) *[]byte {
	key := make([]byte, 16) // NOTE using empty key
	w := siphash.New(key)
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func skein512_256Sum(b *[]byte) *[]byte {
	w := skein.NewHash(32)
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func skein512_512Sum(b *[]byte) *[]byte {
	w := skein.NewHash(64)
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func tiger192Sum(b *[]byte) *[]byte {
	w := tiger.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}

func whirlpoolSum(b *[]byte) *[]byte {
	w := whirlpool.New()
	w.Write(*b)
	res := w.Sum(nil)
	return &res
}
