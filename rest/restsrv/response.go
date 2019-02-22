package restsrv

import (
	"fmt"

	"github.com/mediacoin-pro/core/common/hex"
)

type Response struct {
	Results    interface{} `json:"results,omitempty"`
	NextOffset string      `json:"next_offset,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func NewResponse(res interface{}, nextOffset interface{}, err error) *Response {
	if err != nil {
		return &Response{Error: err.Error()}
	}
	r := &Response{Results: res}
	switch v := nextOffset.(type) {
	case uint64:
		r.NextOffset = "0x" + hex.EncodeUint(v)
	default:
		r.NextOffset = fmt.Sprint(nextOffset)
	}
	return r
}
