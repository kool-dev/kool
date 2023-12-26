package api

import (
	"net/http"
)

func mockHTTPRequester(status int, body string, err error) {
	httpRequester = &fakeHTTP{resp: &http.Response{StatusCode: status, Body: &fakeIOReaderCloser{
		fakeIOReader: fakeIOReader{data: []byte(body), err: err},
	}}}
}
