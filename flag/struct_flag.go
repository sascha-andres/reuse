package flag

import (
	"fmt"
	"reflect"
	"strings"
)

// tagName is the name of the tag for structs
const tagName = "flag"

// Container holds the variables
type Container struct {
	prefix string
	target any

	intVariables    map[string]*int64
	stringVariables map[string]*string
	boolVariables   map[string]*bool
	floatVariables  map[string]*float64
	uintVariables   map[string]*uint64
}

// AddFlagsForStruct adds flags for the given struct.
// The prefix is prepended to the flag name.
// The struct must be a pointer to a struct.
func AddFlagsForStruct(prefix string, s any) (*Container, error) {
	in := reflect.TypeOf(s)

	value := reflect.ValueOf(in)

	if value.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("not a pointer: %s", in.Kind())
	}
	elem := value.Elem()
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: %s", elem.Kind())
	}
	t := elem.Type()

	c := &Container{
		prefix:          prefix,
		target:          s,
		intVariables:    make(map[string]*int64),
		stringVariables: make(map[string]*string),
		boolVariables:   make(map[string]*bool),
		floatVariables:  make(map[string]*float64),
		uintVariables:   make(map[string]*uint64),
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		ef := elem.Field(i)
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
		switch {
		case isInt(field):
			c.intVariables[tag] = Int64(name, 0, usage)
		case field.Type.Kind() == reflect.String:
			c.stringVariables[tag] = String(name, "", usage)
		case field.Type.Kind() == reflect.Bool:
			c.boolVariables[tag] = Bool(name, false, usage)
		case isFloat(field):
			c.floatVariables[tag] = Float64(name, 0.0, usage)
		case isUInt(field):
			c.uintVariables[tag] = Uint64(name, 0, usage)
		case field.Type.Kind() == reflect.Struct:
			for ef.CanAddr() {
				ef = ef.Addr()
			}
			_, err := AddFlagsForStruct(name, ef.Interface())
			if err != nil {
				return nil, err
			}
		default:
			//fmt.Printf("%#v\n", field.Type.Kind().String())
		}
	}

	return c, nil
}

// isInt checks if the given struct field is of type int or int64 or int32 or int16 or int8.
func isInt(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Int ||
		field.Type.Kind() == reflect.Int64 ||
		field.Type.Kind() == reflect.Int32 ||
		field.Type.Kind() == reflect.Int16 ||
		field.Type.Kind() == reflect.Int8
}

// isFloat checks if the given struct field is of type float64 or float32.
func isFloat(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Float64 ||
		field.Type.Kind() == reflect.Float32
}

// isUInt checks if the given struct field is of type uint or uint64 or uint32 or uint16 or uint8.
func isUInt(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Uint ||
		field.Type.Kind() == reflect.Uint64 ||
		field.Type.Kind() == reflect.Uint32 ||
		field.Type.Kind() == reflect.Uint16 ||
		field.Type.Kind() == reflect.Uint8
}

// Parse reads the values from the flags and returns the struct
func (c *Container) Parse() any {
	t := reflect.TypeOf(c.target)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		if strings.Contains(tag, ",") {
			d := strings.SplitN(tag, ",", 2)
			tag = d[0]
		}
		for key, value := range c.stringVariables {
			if key == tag {
				reflect.ValueOf(c.target).Elem().Field(i).SetString(*value)
			}
		}
		for key, value := range c.intVariables {
			if key == tag {
				reflect.ValueOf(c.target).Elem().Field(i).SetInt(*value)
			}
		}
		for key, value := range c.boolVariables {
			if key == tag {
				reflect.ValueOf(c.target).Elem().Field(i).SetBool(*value)
			}
		}
		for key, value := range c.floatVariables {
			if key == tag {
				reflect.ValueOf(c.target).Elem().Field(i).SetFloat(*value)
			}
		}
		for key, value := range c.uintVariables {
			if key == tag {
				reflect.ValueOf(c.target).Elem().Field(i).SetUint(*value)
			}
		}
	}
	return c.target
}
