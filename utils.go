package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
)

func body(x interface{}) ([]byte, bool) {

	switch b := x.(type) {
	case []byte:
		return b, false
	case string:
		return []byte(b), false
	case url.Values:
		return []byte(b.Encode()), false
	case io.ReadCloser:
		defer b.Close()
		bx, err := ioutil.ReadAll(b)
		if err != nil {
			panic(err.Error())
		}
		return bx, true
	case io.Reader:
		bx, err := ioutil.ReadAll(b)
		if err != nil {
			panic(err.Error())
		}
		return bx, true

	default:
		bx, err := json.MarshalIndent(b, "", "    ")

		if err != nil {
			panic(err.Error())
		}
		return bx, false
	}
}
