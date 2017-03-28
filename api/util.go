package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type response struct {
	Result interface{} `json:"result"`
}

// decodeReq decodes request's body to given interface
func decodeReq(r *http.Request, to interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(to); err != nil {
		if err != io.EOF {
			return errors.WithStack(err)
		}
	}
	return nil
}

func queryValue(k string, r *http.Request) string {
	values := r.URL.Query()[k]

	if len(values) != 0 {
		return values[0]
	}

	return ""
}

func queryValueInt(k string, r *http.Request) (int, error) {
	qv := queryValue(k, r)
	if qv == "" {
		return 0, nil
	}
	return strconv.Atoi(qv)
}
