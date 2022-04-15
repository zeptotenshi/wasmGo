//+build tinygo wasm,js

package web

import (
	"fmt"
	"syscall/js"
)

const (
	ELEMENT__tag        = "tagName"
	ELEMENT__id         = "id"
	ELEMENT__class      = "className"
	ELEMENT__innerText  = "innerText"
	element__parent     = "parent"
	element__parentNode = "parentNode"

	function__emit             = "emit"
	function__addEventListener = "addEventListener"
	function__getAttribute     = "getAttribute"
	function__setAttibute      = "setAttribute"
	function__removeAttribute  = "removeAttribute"
	FUNCTION__appendChild      = "appendChild"
	function__removeChild      = "removeChild"

	PROPERTY__value = "value"

	FUNCTION__input_onfocus = "onfocus"
)

func ValidJSValue(_name string, _v js.Value) error {
	if _v.IsUndefined() {
		return fmt.Errorf("js.Value[%s] - undefined", _name)
	}
	if _v.IsNull() {
		return fmt.Errorf("js.Value[%s] - null", _name)
	}
	return nil
}

type Element struct {
	js.Value
	Tag   string
	ID    string
	Class string

	Components map[string]*Component
}

func NewElement(_v js.Value) *Element {
	e := &Element{
		Value:      _v,
		Components: map[string]*Component{},
	}

	id := _v.Get(ELEMENT__id)
	if err := ValidJSValue(ELEMENT__id, id); err == nil {
		e.ID = id.String()
	}
	tag := _v.Get(ELEMENT__tag)
	if err := ValidJSValue(ELEMENT__tag, tag); err == nil {
		e.Tag = tag.String()
	}

	return e
}

// Remove ...
func (elem *Element) Remove() error {
	if elem.Components != nil {
		for k, _ := range elem.Components {
			delete(elem.Components, k)
		}
	}
	parent, err := elem.GetProperty(element__parentNode)
	if err != nil {
		return fmt.Errorf("%s [Remove] [error]: %v", elem, err)
	}
	parent.Call(function__removeChild, elem.Value)
	return nil
}

// SetID ...
func (elem *Element) SetID(_id string) error {
	if err := ValidJSValue(elem.String(), elem.Value); err != nil {
		return fmt.Errorf("%s [SetID] [error]: %v", elem, err)
	}
	elem.ID = _id
	elem.Value.Set(ELEMENT__id, elem.ID)
	return nil
}

// SetClass ...
func (elem *Element) SetClass(_cn string) error {
	if err := ValidJSValue(elem.String(), elem.Value); err != nil {
		return fmt.Errorf("%s [SetClass] [error]: %v", elem, err)
	}
	elem.Value.Set(ELEMENT__class, _cn)
	return nil
}

// SetProperty ...
func (elem *Element) SetProperty(_val interface{}, _names ...string) error {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return fmt.Errorf("%s [SetProperty] [error]: %v", elem, err)
	}

	var pn, pt string
	tv := elem.Value
	for i, n := range _names {
		pt = n

		if i == len(_names)-1 {
			continue
		} else if i == 0 {
			pn = n
		} else {
			pn = fmt.Sprintf("%s.%s", pn, n)
		}

		tv = tv.Get(n)
		if err = ValidJSValue(pn, tv); err != nil {
			return fmt.Errorf("%s [SetProperty] [error]: %v", elem, err)
		}
	}

	tv.Set(pt, _val)
	return nil
}

// SetAttribute ...
func (elem *Element) SetAttribute(_compName string, _vals map[string]interface{}) error {
	if err := ValidJSValue(elem.String(), elem.Value); err != nil {
		return fmt.Errorf("%s [SetAttribute] [error]: %v", elem, err)
	}

	switch len(_vals) {
	case 0:
		elem.Value.Call(function__setAttibute, _compName, "")
	case 1:
		if val, ok := _vals["var"]; ok {
			elem.Value.Call(function__setAttibute, _compName, val)
			return nil
		}
		elem.Value.Call(function__setAttibute, _compName, _vals)
	default:
		elem.Value.Call(function__setAttibute, _compName, _vals)
	}
	return nil
}

// SetAttributes ...
func (elem *Element) SetAttributes(_comps []Component) error {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return fmt.Errorf("%s [SetAttributes] [error]: %v", elem, err)
	}

	var m map[string]interface{}
	for _, v := range _comps {
		if m, err = v.Mapped(); err != nil {
			return fmt.Errorf("%s [SetAttributes] [%s] [error]: %v", elem, v.Name, err)
		}
		elem.SetAttribute(v.Name, m)
	}
	return nil
}

// RemoveAttribute ...
func (elem *Element) RemoveAttribute(_compName string) error {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return fmt.Errorf("%s [RemoveAttribute] [error]: %v", elem, err)
	}
	elem.Value.Call(function__removeAttribute, _compName)
	delete(elem.Components, _compName)
	return nil
}

// SetChild ...
func (elem *Element) SetChild(_v js.Value) error {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return fmt.Errorf("%s [SetChild] [error]: %v", elem, err)
	}
	if err = ValidJSValue("child", _v); err != nil {
		return fmt.Errorf("%s [SetChild] [error]: %v", elem, err)
	}

	elem.Value.Call(FUNCTION__appendChild, _v)
	return nil
}

// RemoveChildByID ...
func (elem *Element) RemoveChildById(_id string) error {
	if err := ValidJSValue(elem.String(), elem.Value); err != nil {
		return fmt.Errorf("%s [RemoveChildById] [error]: %v", elem, err)
	}
	tv := elem.Value.Call(function__getElementById, _id)
	if err := ValidJSValue(_id, tv); err != nil {
		return fmt.Errorf("%s [RemoveChildById] %s [error]: %v", elem, _id, err)
	}
	elem.Value.Call(function__removeChild, tv)
	return nil
}

// StoreCacheValue - NOT IMPLEMENTED
func (elem *Element) StoreCacheValue(_name, _type string, _val interface{}) error {

	return nil
}

// AddEventListener ...
func (elem *Element) AddEventListener(_eventName string, _cb js.Func) error {
	if err := ValidJSValue(elem.String(), elem.Value); err != nil {
		return fmt.Errorf("%s [addEventListener] [error]: %v", elem, err)
	}
	elem.Value.Call(function__addEventListener, _eventName, _cb)
	return nil
}

// Emit ...
func (elem *Element) Emit(_name string, _data map[string]interface{}, _bub bool) {
	elem.Value.Call(function__emit, _name, _data, _bub)
}

// String ...
func (elem *Element) String() string {
	return fmt.Sprintf("[%s]element[%s]", elem.Tag, elem.ID)
}

// Parent ...
func (elem *Element) Parent() string {
	p := elem.Value.Get(element__parent)
	err := ValidJSValue(fmt.Sprintf("%s.%s", elem, element__parent), p)
	if err != nil {
		return ""
	}
	id := p.Get(ELEMENT__id)
	if err = ValidJSValue(fmt.Sprintf("%s.%s.%s", elem, element__parent, ELEMENT__id), id); err != nil {
		return ""
	}
	return id.String()
}

// NOT IMPLEMENTED
func (elem *Element) GetCacheValue(_name string) (interface{}, error) {
	var v interface{}

	return v, nil
}

// GetProperty ...
func (elem *Element) GetProperty(_names ...string) (js.Value, error) {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return js.ValueOf(nil), fmt.Errorf("%s [GetProperty] [error]: %v", elem, err)
	}

	var tv js.Value
	for i, n := range _names {
		if i == 0 {
			tv = elem.Value.Get(n)
		} else {
			tv = tv.Get(n)
		}

		if err = ValidJSValue(n, tv); err != nil {
			return js.ValueOf(nil), fmt.Errorf("%s [GetProperty] [%s] [error]: %v", elem, n, err)
		}
	}
	return tv, nil
}

// GetAttribute ...
func (elem *Element) GetAttribute(_name string) (js.Value, error) {
	err := ValidJSValue(elem.String(), elem.Value)
	if err != nil {
		return js.ValueOf(nil), fmt.Errorf("%s [GetAttribute] [%s] [error]: %v", elem, _name, err)
	}
	v := elem.Value.Call(function__getAttribute, _name)
	if err = ValidJSValue(_name, v); err != nil {
		return js.ValueOf(nil), fmt.Errorf("%s [GetAttribute] [%s] [error]: %v", elem, _name, err)
	}
	return v, nil
}
