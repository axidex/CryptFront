package utils

import "encoding/json"

func PrettyStruct(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "\t")

	return string(s)
}
