package flag

import (
	f "flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// tagName is the name of the tag for structs
const tagName = "flag"

// structNode describes one field discovered while walking a struct for
// AddFlagsForStruct. A node with children == nil is a leaf that has a flag
// registered for it; a node with children != nil is a struct, or
// pointer-to-struct, branch that was recursed into.
type structNode struct {
	index    int
	isPtr    bool
	elemType reflect.Type
	children []*structNode
	kind     reflect.Kind
	fullName string
}

// Container holds the variables
type Container[T any] struct {
	prefix string
	t      *T

	children []*structNode

	intVariables    map[string]*int64
	stringVariables map[string]*string
	boolVariables   map[string]*bool
	floatVariables  map[string]*float64
	uintVariables   map[string]*uint64
}

// AddFlagsForStruct adds flags for the given struct.
// The prefix is prepended to the flag name.
// The struct must be a pointer to a struct.
//
// A field of struct type, or of pointer-to-struct type, is recursed into:
// its own tag becomes the prefix for the fields nested inside it, to
// arbitrary depth. A pointer-to-struct field is left nil by Parse unless at
// least one flag or environment variable inside it was actually supplied.
func AddFlagsForStruct[T any](prefix string, s *T) (*Container[T], error) {
	t := reflect.TypeFor[T]()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: %s", t.Kind())
	}

	c := &Container[T]{
		prefix:          prefix,
		t:               s,
		intVariables:    make(map[string]*int64),
		stringVariables: make(map[string]*string),
		boolVariables:   make(map[string]*bool),
		floatVariables:  make(map[string]*float64),
		uintVariables:   make(map[string]*uint64),
	}
	c.children = buildNodes(prefix, t, c)

	return c, nil
}

// buildNodes walks the fields of t, registers a flag for each tagged leaf
// field (storing the created pointer in c under the field's full flag
// name), and recurses into struct/pointer-to-struct fields.
func buildNodes[T any](prefix string, t reflect.Type, c *Container[T]) []*structNode {
	var nodes []*structNode
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		usage := ""
		if strings.Contains(tag, ",") {
			d := strings.SplitN(tag, ",", 2)
			tag = d[0]
			usage = d[1]
		}
		name := fmt.Sprintf("%v-%v", prefix, tag)

		ft := field.Type
		isPtr := false
		if ft.Kind() == reflect.Pointer {
			isPtr = true
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			nodes = append(nodes, &structNode{
				index:    i,
				isPtr:    isPtr,
				elemType: ft,
				fullName: name,
				children: buildNodes(name, ft, c),
			})
			continue
		}

		if isPtr {
			// pointers to non-struct leaf fields are not supported
			continue
		}

		switch {
		case isIntKind(ft.Kind()):
			c.intVariables[name] = Int64(name, 0, usage)
		case ft.Kind() == reflect.String:
			c.stringVariables[name] = String(name, "", usage)
		case ft.Kind() == reflect.Bool:
			c.boolVariables[name] = Bool(name, false, usage)
		case isFloatKind(ft.Kind()):
			c.floatVariables[name] = Float64(name, 0.0, usage)
		case isUintKind(ft.Kind()):
			c.uintVariables[name] = Uint64(name, 0, usage)
		default:
			continue
		}
		nodes = append(nodes, &structNode{index: i, kind: ft.Kind(), fullName: name})
	}
	return nodes
}

// isIntKind reports whether k is one of the signed integer kinds.
func isIntKind(k reflect.Kind) bool {
	return k == reflect.Int ||
		k == reflect.Int64 ||
		k == reflect.Int32 ||
		k == reflect.Int16 ||
		k == reflect.Int8
}

// isFloatKind reports whether k is one of the floating point kinds.
func isFloatKind(k reflect.Kind) bool {
	return k == reflect.Float64 || k == reflect.Float32
}

// isUintKind reports whether k is one of the unsigned integer kinds.
func isUintKind(k reflect.Kind) bool {
	return k == reflect.Uint ||
		k == reflect.Uint64 ||
		k == reflect.Uint32 ||
		k == reflect.Uint16 ||
		k == reflect.Uint8
}

// Parse reads the values from the flags (and environment variables) and
// writes them into the struct passed to AddFlagsForStruct, returning it.
//
// Must be called after the package-level Parse (flag.Parse) has run, since
// whether a pointer-to-struct branch is populated or left nil depends in
// part on which flags were explicitly set on the command line.
func (c *Container[T]) Parse() *T {
	setFlags := make(map[string]bool)
	Visit(func(fl *f.Flag) {
		setFlags[fl.Name] = true
	})

	v := reflect.ValueOf(c.t).Elem()
	for _, n := range c.children {
		populateNode(n, v, c, setFlags)
	}
	return c.t
}

// populateNode writes n's value(s) into parent and reports whether n (or,
// for a branch, any field beneath it) was explicitly provided via CLI flag
// or environment variable.
func populateNode[T any](n *structNode, parent reflect.Value, c *Container[T], setFlags map[string]bool) bool {
	if n.children != nil {
		target := parent.Field(n.index)
		if n.isPtr {
			target = reflect.New(n.elemType).Elem()
		}

		anyProvided := false
		for _, child := range n.children {
			if populateNode(child, target, c, setFlags) {
				anyProvided = true
			}
		}

		if n.isPtr && anyProvided {
			ptr := reflect.New(n.elemType)
			ptr.Elem().Set(target)
			parent.Field(n.index).Set(ptr)
		}
		return anyProvided
	}

	provided := setFlags[n.fullName]
	if !provided {
		_, provided = os.LookupEnv(envNameForFlagName(n.fullName))
	}

	switch {
	case isIntKind(n.kind):
		parent.Field(n.index).SetInt(*c.intVariables[n.fullName])
	case n.kind == reflect.String:
		parent.Field(n.index).SetString(*c.stringVariables[n.fullName])
	case n.kind == reflect.Bool:
		parent.Field(n.index).SetBool(*c.boolVariables[n.fullName])
	case isFloatKind(n.kind):
		parent.Field(n.index).SetFloat(*c.floatVariables[n.fullName])
	case isUintKind(n.kind):
		parent.Field(n.index).SetUint(*c.uintVariables[n.fullName])
	}

	return provided
}
