package api

import (
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/starudream/go-lib/codec/json"
)

var _v = validator.New()

func V[T any](w http.ResponseWriter, r *http.Request) (v T, err error) {
	defer func() {
		if err != nil {
			ERR(w, http.StatusBadRequest, err.Error())
		}
	}()

	var bs []byte
	bs, err = io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &v)
	if err != nil {
		return
	}

	err = _v.Struct(v)
	if es, ok := err.(validator.ValidationErrors); ok {
		err = es[0]
	}
	return
}
