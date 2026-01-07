package display

import (
	"fmt"
	"reflect"
	"strings"
)

/* =======================
   JSON DISPLAY
   ======================= */

func DisplayJSON(v any) string {
	var b strings.Builder
	writeJSON(&b, reflect.ValueOf(v), 0)
	return b.String()
}

func writeJSON(b *strings.Builder, v reflect.Value, indent int) {
	if !v.IsValid() {
		b.WriteString("null")
		return
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			b.WriteString("null")
			return
		}
		writeJSON(b, v.Elem(), indent)
		return
	}

	switch v.Kind() {

	case reflect.Struct:
		b.WriteString("{\n")
		t := v.Type()
		first := true

		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}

			if !first {
				b.WriteString(",\n")
			}
			first = false

			writeIndent(b, indent+1)
			b.WriteString(fmt.Sprintf(`"%s": `, f.Name))
			writeJSON(b, v.Field(i), indent+1)
		}

		b.WriteString("\n")
		writeIndent(b, indent)
		b.WriteString("}")

	case reflect.Slice, reflect.Array:
		b.WriteString("[\n")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteString(",\n")
			}
			writeIndent(b, indent+1)
			writeJSON(b, v.Index(i), indent+1)
		}
		b.WriteString("\n")
		writeIndent(b, indent)
		b.WriteString("]")

	case reflect.Map:
		b.WriteString("{\n")
		keys := v.MapKeys()
		for i, k := range keys {
			if i > 0 {
				b.WriteString(",\n")
			}
			writeIndent(b, indent+1)
			b.WriteString(fmt.Sprintf(`"%v": `, k.Interface()))
			writeJSON(b, v.MapIndex(k), indent+1)
		}
		b.WriteString("\n")
		writeIndent(b, indent)
		b.WriteString("}")

	case reflect.String:
		b.WriteString(fmt.Sprintf(`"%s"`, v.String()))

	case reflect.Bool:
		b.WriteString(fmt.Sprintf("%v", v.Bool()))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(fmt.Sprintf("%d", v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b.WriteString(fmt.Sprintf("%d", v.Uint()))

	case reflect.Float32, reflect.Float64:
		b.WriteString(fmt.Sprintf("%v", v.Float()))

	default:
		b.WriteString(fmt.Sprintf(`"%v"`, v.Interface()))
	}
}

/* =======================
   STRUCT DISPLAY
   ======================= */

func DisplayStruct(v any) {
	var b strings.Builder
	writeStruct(&b, reflect.ValueOf(v), 0)
	print(b.String())
}

func writeStruct(b *strings.Builder, v reflect.Value, indent int) {
	if !v.IsValid() {
		b.WriteString("nil")
		return
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			b.WriteString("nil")
			return
		}
		writeStruct(b, v.Elem(), indent)
		return
	}

	switch v.Kind() {

	case reflect.Struct:
		t := v.Type()
		b.WriteString(t.Name() + " {\n")

		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}

			writeIndent(b, indent+1)
			b.WriteString(f.Name + ": ")
			writeStruct(b, v.Field(i), indent+1)
			b.WriteString("\n")
		}

		writeIndent(b, indent)
		b.WriteString("}")

	case reflect.Slice, reflect.Array:
		b.WriteString("[\n")
		for i := 0; i < v.Len(); i++ {
			writeIndent(b, indent+1)
			writeStruct(b, v.Index(i), indent+1)
			b.WriteString("\n")
		}
		writeIndent(b, indent)
		b.WriteString("]")

	case reflect.Map:
		b.WriteString("{\n")
		for _, k := range v.MapKeys() {
			writeIndent(b, indent+1)
			b.WriteString(fmt.Sprintf("%v: ", k.Interface()))
			writeStruct(b, v.MapIndex(k), indent+1)
			b.WriteString("\n")
		}
		writeIndent(b, indent)
		b.WriteString("}")

	default:
		b.WriteString(fmt.Sprintf("%v", v.Interface()))
	}
}

func writeIndent(b *strings.Builder, level int) {
	b.WriteString(strings.Repeat("  ", level))
}
