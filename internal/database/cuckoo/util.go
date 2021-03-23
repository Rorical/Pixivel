package cuckoo

import (
	"bytes"
	"encoding/binary"
	"sync/atomic"

	metro "github.com/dgryski/go-metro"
)

var (
	altHash = [256]uint{}
	masks   = [65]uint{}
)

func init() {
	for i := 0; i < 256; i++ {
		altHash[i] = (uint(metro.Hash64([]byte{byte(i)}, 1337)))
	}
	for i := uint(0); i <= 64; i++ {
		masks[i] = (1 << i) - 1
	}
}

func getAltIndex(fp byte, i uint, bucketPow uint) uint {
	mask := masks[bucketPow]
	hash := altHash[fp] & mask
	return (i & mask) ^ hash
}

func getFingerprint(hash uint64) byte {
	// Use least significant bits for fingerprint.
	fp := byte(hash%255 + 1)
	return fp
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, bucketPow uint) (uint, uint, byte) {
	hash := metro.Hash64(data, 1337)
	f := getFingerprint(hash)
	// Use most significant bits for deriving index.
	i1 := uint(hash>>32) & masks[bucketPow]
	i2 := getAltIndex(f, i1, bucketPow)
	return i1, i2, f
}

func getNextPow2(n uint64) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return uint(n)
}

const a = uint64(2862933555777941757)
const b = uint64(3037000493)

// Linear Congruential Generator
// See https://link.springer.com/chapter/10.1007/978-1-4615-2317-8_3
type LCG struct {
	state uint64
}

func (l *LCG) Intn(n int) int {
	for {
		state := atomic.LoadUint64(&l.state)

		newState := a*state + b

		// Replace only if the state is still the same
		swapped := atomic.CompareAndSwapUint64(&l.state, state, newState)
		if swapped {
			return int(newState % uint64(n))
		}
	}
}

func UintIn(n uint) []byte {
	data := uint64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func UintOut(bye []byte) uint {
	bytebuff := bytes.NewBuffer(bye)
	var data uint64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return uint(data)
}
