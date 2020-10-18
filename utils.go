package pixivel

import (
	"crypto/md5"
	"encoding/json"
)

//Hash will md5 the struct
func HashStruct(item interface{}) [16]byte {
	jsonBytes, _ := json.Marshal(item)
	return md5.Sum(jsonBytes)
}
