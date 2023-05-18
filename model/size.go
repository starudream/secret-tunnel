package model

import (
	"github.com/starudream/secret-tunnel/internal/unitx"
)

type Size uint

func (x Size) String() string {
	return unitx.HumanSize(float64(x))
}
