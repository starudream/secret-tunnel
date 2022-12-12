package tablew

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/olekukonko/tablewriter"
)

func SetStructs(t *tablewriter.Table, v any) error {
	if v == nil {
		return errors.New("nil value")
	}
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	switch vt.Kind() {
	case reflect.Slice, reflect.Array:
		if vv.Len() < 1 {
			return errors.New("empty value")
		}

		// check first element to set header
		first := vv.Index(0)
		e := first.Type()
		switch e.Kind() {
		case reflect.Struct:
			// OK
		case reflect.Ptr:
			if first.IsNil() {
				return errors.New("the first element is nil")
			}
			e = first.Elem().Type()
			if e.Kind() != reflect.Struct {
				return fmt.Errorf("invalid kind %s", e.Kind())
			}
		default:
			return fmt.Errorf("invalid kind %s", e.Kind())
		}
		n := e.NumField()
		headers := make([]string, n)
		for i := 0; i < n; i++ {
			f := e.Field(i)
			header := f.Tag.Get("tablewriter")
			if header == "" {
				header = f.Name
			}
			headers[i] = header
		}
		t.SetHeader(headers)

		for i := 0; i < vv.Len(); i++ {
			item := reflect.Indirect(vv.Index(i))
			itemType := reflect.TypeOf(item)
			switch itemType.Kind() {
			case reflect.Struct:
				// OK
			default:
				return fmt.Errorf("invalid item type %v", itemType.Kind())
			}
			if !item.IsValid() {
				// skip rendering
				continue
			}
			nf := item.NumField()
			if n != nf {
				return errors.New("invalid num of field")
			}
			rows := make([]string, nf)
			for j := 0; j < nf; j++ {
				f := reflect.Indirect(item.Field(j))
				if f.Kind() == reflect.Ptr {
					f = f.Elem()
				}
				if f.IsValid() {
					switch s := f.Interface().(type) {
					case time.Time:
						if s.IsZero() {
							rows[j] = ""
						} else {
							rows[j] = s.Format(time.RFC3339)
						}
					case fmt.Stringer:
						rows[j] = s.String()
					default:
						rows[j] = fmt.Sprint(f)
					}
				} else {
					rows[j] = "nil"
				}
			}
			t.Append(rows)
		}
	default:
		return fmt.Errorf("invalid type %T", v)
	}
	return nil
}

func PrintStruct(v any) string {
	ts := reflect.SliceOf(reflect.TypeOf(v))
	ss := reflect.MakeSlice(ts, 1, 1)
	vs := reflect.New(ss.Type())
	vs.Elem().Set(ss)
	vs.Elem().Index(0).Set(reflect.ValueOf(v))
	return PrintStructs(vs.Elem().Interface())
}

func PrintStructs(v any) string {
	bb := &bytes.Buffer{}
	w := tablewriter.NewWriter(bb)
	err := SetStructs(w, v)
	if err != nil {
		return ""
	}
	w.Render()
	return bb.String()
}
