//+build tinygo wasm,js

package web

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	GOWEB_ERROR_COMPONENT_VALS_NIL     = errors.New("component vals map nil")
	GOWEB_ERROR_INVALID_ATTRIBUTE_TYPE = errors.New("invalid component attribute type")
)

type Attribute struct {
	Type  string
	Value string
}

func (a Attribute) String() string {
	return fmt.Sprintf("type[%s] value[%s]", a.Type, a.Value)
}

type Component struct {
	Name string
	Vals map[string]Attribute
}

func (c *Component) Mapped() (map[string]interface{}, error) {
	if c.Vals == nil {
		return nil, GOWEB_ERROR_COMPONENT_VALS_NIL
	}

	r := map[string]interface{}{}

	if len(c.Vals) == 0 {
		return r, nil
	}

	for k, v := range c.Vals {
		switch v.Type {
		case "number":
			f, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				return nil, err
			}
			r[k] = f

		case "bool":
			b, err := strconv.ParseBool(v.Value)
			if err != nil {
				return nil, err
			}
			r[k] = b

		case "string":
			r[k] = v.Value

		default:
			return r, GOWEB_ERROR_INVALID_ATTRIBUTE_TYPE
		}
	}

	return r, nil
}

func (c *Component) String() string {
	s := fmt.Sprintf("component[%s]", c.Name)
	for i, v := range c.Vals {
		s = fmt.Sprintf("%s\n\t\t%s:\n\t\t%s", s, i, v.String())
	}
	return s
}
