package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// generateReflect generates a self-contained Go program that extracts component
// metadata via reflection. The entire reflection logic is inlined.
func generateReflect(mainCompName string, childCompNames []string) string {
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
`)

	return sb.String()
}

// parseReflectOutput parses the output from reflect.go into component structs
func parseReflectOutput(output string) ([]*component, []staticRoute, error) {
	var components []*component
	var routes []staticRoute
	var current *component

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

	return components, routes, nil
}

// staticRoute represents a route for static HTML generation
type staticRoute struct {
	path     string
	htmlFile string
}
