package config

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

//Hash will md5 the struct
func HashStruct(item interface{}) string {
	jsonBytes, _ := json.Marshal(item)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}
