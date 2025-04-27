package services

import "net/http"

type ProductMsgProc interface {
	Decode(data []byte) (interface{}, error)
	Validate(v interface{}) error
	ProcessMsg(v interface{}, r *http.Request) (interface{}, error)
	Encode(v interface{}) ([]byte, int, error)
}
