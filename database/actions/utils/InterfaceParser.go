package utils

import (
	"bytes"
	"encoding/gob"
)

// ---------------------------------
//          DEPRECATED
// This function is no longer used
// ----------------------------------
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
