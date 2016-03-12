package gohash

var (
	algos = map[string]int{
		"adler32":      32,
		"blake224":     224,
		"blake256":     256,
		"blake384":     384,
		"blake512":     512,
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
		//"skein256-256": 256,
		"skein512-256": 256,
		"skein512-512": 512,
		"tiger192":     192,
		"whirlpool":    512,
	}
)
