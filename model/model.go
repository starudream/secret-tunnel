package model

import (
	"gorm.io/gorm"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
	"github.com/starudream/go-lib/sqlite/v2"

	"github.com/starudream/secret-tunnel/util"
)

var Expr = gorm.Expr

var tables = []any{
	&Client{},
	&Task{},
}

func init() {
	osutil.PanicErr(sqlite.DB().AutoMigrate(tables...))

	gh.Silently(UpdateAllClientOffline())
}

type Size uint

func (s Size) String() string {
	return util.IBytes(uint64(s))
}
