package util

import (
	"github.com/google/uuid"

	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

func UUID() string {
	v, err := uuid.NewRandom()
	osutil.PanicErr(err)
	return v.String()
}

func UUIDShort() string {
	s := UUID()
	return s[0:8] + s[9:13] + s[14:18] + s[19:23] + s[24:36]
}
