package utl

import "encoding/json"

func BsInto[T any](bs []byte) (*T, error) {
	out := new(T)
	return out, json.Unmarshal(bs, out)
}

func BsIntoPnk[T any](bs []byte) *T {
	out := new(T)
	if err := json.Unmarshal(bs, out); err != nil {
		panic(err.Error())
	}
	return out
}
