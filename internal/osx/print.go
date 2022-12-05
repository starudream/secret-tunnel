package osx

import (
	"fmt"
	"os"
	"strings"
)

func P(v ...any) {
	if len(v) == 0 {
		os.Exit(0)
	}
	c, s, w := 0, "", os.Stdout
	nh := func() {
		if len(v) >= 2 {
			switch y := v[1].(type) {
			case string:
				v = v[1:]
				s = y
			}
		}
	}
	switch x := v[0].(type) {
	case string:
		s = x
	case error:
		if x != nil {
			c, s, w = 1, x.Error(), os.Stderr
			v = v[:1]
		} else {
			nh()
		}
	case nil:
		nh()
	default:
		c, s = 1, fmt.Sprint(x)
	}
	if len(v) >= 2 {
		s = fmt.Sprintf(s, v[1:]...)
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	if s != "" {
		_, _ = fmt.Fprintf(w, s)
	}
	os.Exit(c)
}
