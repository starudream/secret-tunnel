package util

import (
	"fmt"

	"github.com/starudream/go-lib/tablew/v2"
)

func TablePrint(vs any) error {
	fmt.Println(tablew.Structs(vs, func(w *tablew.Table) { w.SetColWidth(40) }))
	return nil
}
