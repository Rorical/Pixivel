package cuckoo

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const maxCuckooCount = 500

// Filter is a probabilistic counter
type Filter struct {
	buckets   []bucket
	count     uint
	bucketPow uint
	FilePath  string
	gen       LCG
}

// NewFilter returns a new cuckoofilter with a given capacity.
// A capacity of 1000000 is a normal default, which allocates
// about ~1MB on 64-bit machines.
func NewFilter(capacity uint, path string) *Filter {
	cf, err := ReadFile(path)
	if err != nil {
		capacity = getNextPow2(uint64(capacity)) / bucketSize
		if capacity == 0 {
			capacity = 1
		}
		buckets := make([]bucket, capacity)
		cf = &Filter{
			buckets:   buckets,
			count:     0,
			bucketPow: uint(bits.TrailingZeros(capacity)),
			FilePath:  path,
		}
	}

	return cf
}

// Lookup returns true if data is in the counter
func (cf *Filter) Lookup(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bucketPow)
	b1, b2 := cf.buckets[i1], cf.buckets[i2]
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

// Expand expands the buckets when it was fulfilled
func (cf *Filter) Expand() {
	capacity := uint(1) << (uint(cf.bucketPow + 1))
	buckets := make([]bucket, capacity)
	copy(buckets, cf.buckets)
	cf.buckets = buckets
	cf.bucketPow = uint(bits.TrailingZeros(capacity))
}

func (cf *Filter) Reset() {
	for i := range cf.buckets {
		cf.buckets[i].reset()
	}
	cf.count = 0
}

func (cf *Filter) randi(i1, i2 uint) uint {
	if cf.gen.Intn(2) == 0 {
		return i1
	}
	return i2
}

// Insert inserts data into the counter and returns true upon success
func (cf *Filter) Insert(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bucketPow)
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, cf.randi(i1, i2))
}

// InsertUnique inserts data into the counter if not exists and returns true upon success
func (cf *Filter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

func (cf *Filter) insert(fp byte, i uint) bool {
	if cf.buckets[i].insert(fp) {
		cf.count++
		return true
	}
	return false
}

func (cf *Filter) reinsert(fp byte, i uint) bool {
	count := 0
reinsert:
	for k := 0; k < maxCuckooCount; k++ {
		j := cf.gen.Intn(bucketSize)
		oldfp := fp
		fp = cf.buckets[i][j]
		cf.buckets[i][j] = oldfp

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, cf.bucketPow)
		if cf.insert(fp, i) {
			return true
		}
	}
	cf.Expand()
	count++
	if count < 3 {
		goto reinsert
	}
	return false
}

// Delete data from counter if exists and return if deleted or not
func (cf *Filter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bucketPow)
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

func (cf *Filter) delete(fp byte, i uint) bool {
	if cf.buckets[i].delete(fp) {
		cf.count--
		return true
	}
	return false
}

// Count returns the number of items in the counter
func (cf *Filter) Count() uint {
	return cf.count
}

// Encode returns a byte slice representing a Cuckoofilter
func (cf *Filter) Encode() []byte {
	bytes := make([]byte, len(cf.buckets)*bucketSize)
	for i, b := range cf.buckets {
		for j, f := range b {
			index := (i * len(b)) + j
			bytes[index] = f
		}
	}
	return bytes
}

// Decode returns a Cuckoofilter from a byte slice
func Decode(bytes []byte) (*Filter, error) {
	var count uint
	if len(bytes)%bucketSize != 0 {
		return nil, fmt.Errorf("expected bytes to be multiple of %d, got %d", bucketSize, len(bytes))
	}
	buckets := make([]bucket, len(bytes)/4)
	for i, b := range buckets {
		for j := range b {
			index := (i * len(b)) + j
			if bytes[index] != 0 {
				buckets[i][j] = bytes[index]
				count++
			}
		}
	}
	return &Filter{
		buckets:   buckets,
		count:     count,
		bucketPow: uint(bits.TrailingZeros(uint(len(buckets)))),
	}, nil
}

func (cf *Filter) SaveFile() error {
	file, err := os.Create(cf.FilePath)
	if err != nil {
		return err
	}
	_, err = file.Write(cf.Encode())
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
func ReadFile(path string) (*Filter, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return nil, err
	}
	cf, err := Decode(bytes)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return cf, nil
}

func (cf *Filter) ReadFile() error {
	file, err := os.Open(cf.FilePath)
	if err != nil {
		return err
	}
	stats, err := file.Stat()
	if err != nil {
		return err
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return err
	}
	cf, err = Decode(bytes)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
