package atcoder

import (
	"context"
	"net/http"
)

func NewAtCoderTester(rt http.RoundTripper) *AtCoder {
	ac := NewAtCoder(context.Background(), "")
	ac.client.client.Transport = rt
	return ac
}
