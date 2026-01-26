package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// generateReflect generates a self-contained Go program that extracts component
// metadata via reflection. The entire reflection logic is inlined.
func generateReflect(mainCompName string, childCompNames []string, codecTypes []string) string {
	var sb strings.Builder

	sb.WriteString(`//go:build !wasm

package main

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
)

func main() {
`)

	// Main component first
	fmt.Fprintf(&sb, "\textract(&%s{}, true)\n", mainCompName)

	// Child components
	for _, name := range childCompNames {
		fmt.Fprintf(&sb, "\textract(&%s{}, false)\n", name)
	}

	// Codec types
	for _, name := range codecTypes {
		fmt.Fprintf(&sb, "\textractCodec(&%s{})\n", name)
	}

	sb.WriteString(`}

func extract(v interface{}, isMain bool) {
	rv := reflect.ValueOf(v)
	rt := rv.Type().Elem()

	fmt.Printf("COMPONENT:%s\n", rt.Name())

	if m := rv.MethodByName("Template"); m.IsValid() {
		if r := m.Call(nil); len(r) > 0 {
			fmt.Printf("TEMPLATE:%s\n", base64.StdEncoding.EncodeToString([]byte(r[0].String())))
		}
	}

	if m := rv.MethodByName("Style"); m.IsValid() {
		if r := m.Call(nil); len(r) > 0 {
			fmt.Printf("STYLE:%s\n", base64.StdEncoding.EncodeToString([]byte(r[0].String())))
		}
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if st, vt, kt := parseType(f.Type.String()); st != "" {
			fmt.Printf("FIELD:%s:%s:%s:%s\n", f.Name, st, vt, kt)
		}
	}

	for i := 0; i < rv.NumMethod(); i++ {
		fmt.Printf("METHOD:%s\n", rv.Type().Method(i).Name)
	}

	if isMain {
		if m := rv.MethodByName("Routes"); m.IsValid() {
			if r := m.Call(nil); len(r) > 0 {
				routes := r[0]
				for i := 0; i < routes.Len(); i++ {
					route := routes.Index(i)
					fmt.Printf("ROUTES:%s:%s\n", route.FieldByName("Path").String(), route.FieldByName("HTMLFile").String())
				}
			}
		}
	}

	fmt.Println("END")
}

func parseType(t string) (store, val, key string) {
	t = strings.TrimPrefix(t, "*")
	if i := strings.LastIndex(t, "."); i != -1 {
		t = t[i+1:]
	}
	if t == "LocalStore" {
		return "LocalStore", "string", ""
	}
	if i := strings.Index(t, "["); i != -1 {
		store = t[:i]
		params := t[i+1 : len(t)-1]
		switch store {
		case "Store", "List":
			return store, params, ""
		case "Map":
			d := 0
			for j, c := range params {
				switch c {
				case '[':
					d++
				case ']':
					d--
				case ',':
					if d == 0 {
						return store, strings.TrimSpace(params[j+1:]), strings.TrimSpace(params[:j])
					}
				}
			}
		}
	}
	return "", "", ""
}

func extractCodec(v interface{}) {
	rt := reflect.TypeOf(v).Elem()
	fmt.Printf("CODEC:%s\n", rt.Name())
	extractCodecFields(rt, "")
	fmt.Println("CODEC_END")
}

func extractCodecFields(rt reflect.Type, prefix string) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		jsTag := f.Tag.Get("js")
		if jsTag == "-" {
			continue
		}
		if jsTag == "" {
			jsTag = f.Name
		}
		// Handle omitempty
		omit := false
		if idx := strings.Index(jsTag, ","); idx != -1 {
			if strings.Contains(jsTag[idx:], "omitempty") {
				omit = true
			}
			jsTag = jsTag[:idx]
		}

		fieldPath := f.Name
		if prefix != "" {
			fieldPath = prefix + "." + f.Name
		}

		ft := f.Type
		// Dereference pointer
		isPtr := false
		if ft.Kind() == reflect.Ptr {
			isPtr = true
			ft = ft.Elem()
		}

		kind := ft.Kind().String()
		typeName := ft.String()

		// Handle anonymous/embedded structs
		if ft.Kind() == reflect.Struct && f.Anonymous {
			extractCodecFields(ft, fieldPath)
			continue
		}

		// Handle named nested structs
		if ft.Kind() == reflect.Struct {
			fmt.Printf("CODEC_FIELD:%s:%s:%s:struct:%v:%v\n", fieldPath, jsTag, typeName, isPtr, omit)
			extractCodecFields(ft, fieldPath)
			continue
		}

		// Handle slices
		if ft.Kind() == reflect.Slice {
			elemType := ft.Elem()
			elemKind := elemType.Kind().String()
			if elemType.Kind() == reflect.Struct {
				fmt.Printf("CODEC_FIELD:%s:%s:%s:slice_struct:%v:%v\n", fieldPath, jsTag, elemType.String(), isPtr, omit)
				extractCodecFields(elemType, fieldPath+"[]")
			} else {
				fmt.Printf("CODEC_FIELD:%s:%s:%s:slice_%s:%v:%v\n", fieldPath, jsTag, elemType.String(), elemKind, isPtr, omit)
			}
			continue
		}

		// Handle maps
		if ft.Kind() == reflect.Map {
			keyType := ft.Key().String()
			valType := ft.Elem()
			valKind := valType.Kind().String()
			if valType.Kind() == reflect.Struct {
				fmt.Printf("CODEC_FIELD:%s:%s:%s:map_%s_struct:%v:%v\n", fieldPath, jsTag, keyType, valKind, isPtr, omit)
				extractCodecFields(valType, fieldPath+"[]")
			} else {
				fmt.Printf("CODEC_FIELD:%s:%s:%s:map_%s_%s:%v:%v\n", fieldPath, jsTag, keyType, valKind, valType.String(), isPtr, omit)
			}
			continue
		}

		fmt.Printf("CODEC_FIELD:%s:%s:%s:%s:%v:%v\n", fieldPath, jsTag, typeName, kind, isPtr, omit)
	}
}
`)

	return sb.String()
}

// codecType represents a type that needs FromJS/ToJS codegen
type codecType struct {
	name   string
	fields []codecField
}

// codecField represents a field in a codec type
type codecField struct {
	path     string // e.g., "RAW.PRICE" for nested
	jsName   string // from js tag or field name
	typeName string // Go type string
	kind     string // "string", "int", "struct", "slice_struct", etc.
	isPtr    bool
	omit     bool // omitempty
}

// parseReflectOutput parses the output from reflect.go into component structs
func parseReflectOutput(output string) ([]*component, []staticRoute, []*codecType, error) {
	var components []*component
	var routes []staticRoute
	var codecTypes []*codecType
	var current *component
	var currentCodec *codecType

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "COMPONENT":
			current = &component{name: value}
			components = append(components, current)

		case "TEMPLATE":
			if current != nil {
				if decoded, err := base64.StdEncoding.DecodeString(value); err == nil {
					current.template = string(decoded)
				}
			}

		case "STYLE":
			if current != nil {
				if decoded, err := base64.StdEncoding.DecodeString(value); err == nil {
					current.style = string(decoded)
				}
			}

		case "CODEC":
			currentCodec = &codecType{name: value}
			codecTypes = append(codecTypes, currentCodec)

		case "CODEC_FIELD":
			if currentCodec != nil {
				// Format: path:jsName:typeName:kind:isPtr:omit
				fieldParts := strings.Split(value, ":")
				if len(fieldParts) >= 6 {
					currentCodec.fields = append(currentCodec.fields, codecField{
						path:     fieldParts[0],
						jsName:   fieldParts[1],
						typeName: fieldParts[2],
						kind:     fieldParts[3],
						isPtr:    fieldParts[4] == "true",
						omit:     fieldParts[5] == "true",
					})
				}
			}

		case "CODEC_END":
			currentCodec = nil

		case "FIELD":
			if current != nil {
				fieldParts := strings.Split(value, ":")
				if len(fieldParts) >= 3 {
					field := storeField{
						name:      fieldParts[0],
						storeType: fieldParts[1],
						valueType: fieldParts[2],
					}
					if len(fieldParts) >= 4 {
						field.keyType = fieldParts[3]
					}
					current.fields = append(current.fields, field)
				}
			}

		case "METHOD":
			if current != nil {
				current.methods = append(current.methods, value)
				if value == "OnMount" {
					current.hasOnMount = true
				}
				if value == "OnCreate" {
					current.hasOnCreate = true
				}
				if value == "OnUnmount" {
					current.hasOnUnmount = true
				}
			}

		case "ROUTES":
			routeParts := strings.SplitN(value, ":", 2)
			if len(routeParts) == 2 {
				routes = append(routes, staticRoute{
					path:     routeParts[0],
					htmlFile: routeParts[1],
				})
			}
		}
	}

	return components, routes, codecTypes, nil
}

// staticRoute represents a route for static HTML generation
type staticRoute struct {
	path     string
	htmlFile string
}
