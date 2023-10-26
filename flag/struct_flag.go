package flag

import (
	"fmt"
	"reflect"
	"strings"
)

// tagName is the name of the tag for structs
const tagName = "flag"

// Container holds the variables
type Container[T any] struct {
	prefix string
	t      T

	intVariables    map[string]*int64
	stringVariables map[string]*string
	boolVariables   map[string]*bool
	floatVariables  map[string]*float64
	uintVariables   map[string]*uint64
}

// AddFlagsForStruct adds flags for the given struct.
// The prefix is prepended to the flag name.
// The struct must be a pointer to a struct.
func AddFlagsForStruct[T any](prefix string, s T) (*Container[T], error) {
	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(s)
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

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
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
		var name = fmt.Sprintf("%v-%v", prefix, tag)
		if field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Int64 || field.Type.Kind() == reflect.Int32 || field.Type.Kind() == reflect.Int16 || field.Type.Kind() == reflect.Int8 {
			c.intVariables[tag] = Int64(name, reflect.ValueOf(&c.t).Elem().Field(i).Int(), usage)
		}
		if field.Type.Kind() == reflect.String {
			c.stringVariables[tag] = String(name, reflect.ValueOf(&c.t).Elem().Field(i).String(), usage)
		}
		if field.Type.Kind() == reflect.Bool {
			c.boolVariables[tag] = Bool(name, reflect.ValueOf(&c.t).Elem().Field(i).Bool(), usage)
		}
		if field.Type.Kind() == reflect.Float64 || field.Type.Kind() == reflect.Float32 {
			c.floatVariables[tag] = Float64(name, reflect.ValueOf(&c.t).Elem().Field(i).Float(), usage)
		}
		if field.Type.Kind() == reflect.Uint || field.Type.Kind() == reflect.Uint64 || field.Type.Kind() == reflect.Uint32 || field.Type.Kind() == reflect.Uint16 || field.Type.Kind() == reflect.Uint8 {
			c.uintVariables[tag] = Uint64(name, reflect.ValueOf(&c.t).Elem().Field(i).Uint(), usage)
		}
	}

	return c, nil
}

func (c *Container[T]) Parse() *T {
	t := reflect.TypeOf(c.t)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		if strings.Contains(tag, ",") {
			tag = strings.SplitN(tag, ",", 2)[0]
		}
		for key, value := range c.stringVariables {
			if key == tag {
				reflect.ValueOf(&c.t).Elem().Field(i).SetString(*value)
			}
		}
		for key, value := range c.intVariables {
			if key == tag {
				reflect.ValueOf(&c.t).Elem().Field(i).SetInt(int64(*value))
			}
		}
		for key, value := range c.boolVariables {
			if key == tag {
				reflect.ValueOf(&c.t).Elem().Field(i).SetBool(*value)
			}
		}
		for key, value := range c.floatVariables {
			if key == tag {
				reflect.ValueOf(&c.t).Elem().Field(i).SetFloat(*value)
			}
		}
		for key, value := range c.uintVariables {
			if key == tag {
				reflect.ValueOf(&c.t).Elem().Field(i).SetUint(*value)
			}
		}
	}
	return &c.t
}
