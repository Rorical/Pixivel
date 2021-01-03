package levelgo

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrNotFound         = leveldb.ErrNotFound
	ErrReadOnly         = leveldb.ErrReadOnly
	ErrSnapshotReleased = leveldb.ErrSnapshotReleased
	ErrIterReleased     = leveldb.ErrIterReleased
	ErrClosed           = leveldb.ErrClosed
)

func ErrorToErrCode(err error) int32 {
	var errorCode int32
	errorCode = -1
	switch err {
	case ErrNotFound:
		errorCode = 0
	case ErrReadOnly:
		errorCode = 1
	case ErrSnapshotReleased:
		errorCode = 2
	case ErrIterReleased:
		errorCode = 3
	case ErrClosed:
		errorCode = 4
	}
	return errorCode
}

func ErrCodeToError(errorCode int32) error {
	switch errorCode {
	case 0:
		return ErrNotFound
	case 1:
		return ErrReadOnly
	case 2:
		return ErrSnapshotReleased
	case 3:
		return ErrIterReleased
	case 4:
		return ErrClosed
	}
	return nil
}
