package pixivel

import (
	"bytes"
	"encoding/binary"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBSetting struct {
	File string
}

type LevelDB struct {
	db       *leveldb.DB
	NotFound error
}

type LevelDBBatchOperations struct {
	batch *leveldb.Batch
}

func GetLevelDB() *LevelDB {
	var err error
	db, err := leveldb.OpenFile(leveldbConf.File, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDB{
		db:       db,
		NotFound: leveldb.ErrNotFound,
	}
}

func (self *LevelDB) Get(key []byte) ([]byte, error) {
	data, err := self.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (self *LevelDB) Set(key []byte, value []byte) error {
	err := self.db.Put(key, value, nil)
	return err
}

func (self *LevelDB) Has(key []byte) (bool, error) {
	has, err := self.db.Has(key, nil)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (self *LevelDB) Del(key []byte) error {
	err := self.db.Delete(key, nil)
	return err
}

func (self *LevelDB) GetBatch() *LevelDBBatchOperations {
	batch := new(leveldb.Batch)
	return &LevelDBBatchOperations{
		batch: batch,
	}
}

func (self *LevelDBBatchOperations) Set(key []byte, value []byte) {
	self.batch.Put(key, value)
}

func (self *LevelDBBatchOperations) Del(key []byte) {
	self.batch.Delete(key)
}

func (self *LevelDB) RunBatch(batch *LevelDBBatchOperations) error {
	err := self.db.Write(batch.batch, nil)
	return err
}

func (self *LevelDB) CloseLevelDB() {
	self.db.Close()
}

//Some utils
func (self *LevelDB) StringOut(bye []byte) string {
	return string(bye)
}

func (self *LevelDB) StringIn(strings string) []byte {
	return []byte(strings)
}

func (self *LevelDB) IntIn(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func (self *LevelDB) IntOut(bye []byte) int {
	bytebuff := bytes.NewBuffer(bye)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func (self *LevelDB) UintIn(n uint) []byte {
	data := uint64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func (self *LevelDB) UintOut(bye []byte) uint {
	bytebuff := bytes.NewBuffer(bye)
	var data uint64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return uint(data)
}
