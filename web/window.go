//+build tinygo wasm,js

package web

import (
	"fmt"
	"strings"
	"syscall/js"
)

const (
	window   = "window"
	document = "document"
	title    = "title"
	worker   = "worker"
	cookie   = "cookie"

	WINDOW__location           = "location"
	FUNCTION__location_replace = "replace"

	FUNCTION__form_onsubmit       = "onsubmit"
	FUNCTION__form_preventDefault = "preventDefault"

	function__createElement  = "createElement"
	function__getElementById = "getElementById"
	function__querySelector  = "querySelector"
	function__remove         = "remove"
)

type Window struct {
	Title   string
	Version string

	value    js.Value
	document js.Value

	*Logger
}

// NewWindow ...
func NewWindow() *Window {
	win := &Window{
		value:    js.Global(),
		document: js.ValueOf(nil),
		Logger:   NewLogger(),
	}

	tv := win.value.Get(document)
	if err := ValidJSValue(fmt.Sprintf("%s.%s", window, document), tv); err != nil {
		win.Logger.Prefix = worker
	} else {
		win.document = tv
		tv = tv.Get(title)
		if err = ValidJSValue(fmt.Sprintf("%s.%s.%s", window, document, title), tv); err == nil {
			win.Title = tv.String()
			win.Logger.Prefix = win.Title
		}
	}

	return win
}

// NewElementWithTag ...
func (w *Window) NewElementWithTag(_tag string) *Element {
	elem := NewElement(js.ValueOf(nil))
	if err := ValidJSValue(document, w.document); err != nil {
		return elem
	}
	elem.Value = w.document.Call(function__createElement, _tag)
	return elem
}

// NewElementWithValue ...
func (w *Window) NewElementWithValue(_v js.Value) *Element {
	return NewElement(_v)
}

// ElementById ...
func (w *Window) ElementById(_id string) *Element {
	tv, err := w.GetValueById(_id)
	if err != nil {
		w.Logger.Error(fmt.Errorf("[window] [ElementByID] [error]: %v", err))
		return nil
	}
	if err = ValidJSValue(_id, tv); err != nil {
		w.Logger.Error(fmt.Errorf("[window] [ElementByID] [error]: %v", err))
		return nil
	}
	return NewElement(tv)
}

// GetGlobal ...
func (w *Window) GetGlobal(_name string) (js.Value, error) {
	err := ValidJSValue(window, w.value)
	if err != nil {
		return js.ValueOf(nil), fmt.Errorf("[window] [GetGlobal] [error]: %v", err)
	}
	tv := w.value.Get(_name)
	if err = ValidJSValue(_name, tv); err != nil {
		return js.ValueOf(nil), fmt.Errorf("[window] [GetGlobal] [error]: %v", _name, err)
	}
	return tv, nil
}

// SetGlobal ...
func (w *Window) SetGlobal(_name string, _val interface{}) error {
	err := ValidJSValue(window, w.value)
	if err != nil {
		return fmt.Errorf("[window] [SetGlobal] [error]: %v", err)
	}
	w.value.Set(_name, _val)
	return nil
}

// GetValueById ...
func (w *Window) GetValueById(_id string) (js.Value, error) {
	if err := ValidJSValue(document, w.document); err != nil {
		return js.ValueOf(nil), fmt.Errorf("[window] [GetValueById] [error]: %v", err)
	}
	return w.document.Call(function__getElementById, _id), nil
}

// GetCookie ...
func (w *Window) GetCookie(_name string) (string, error) {
	err := ValidJSValue(document, w.document)
	if err != nil {
		return "", fmt.Errorf("[window] [GetCookie] [error]: %v", err)
	}
	tv := w.document.Get(cookie)
	if err = ValidJSValue(cookie, tv); err != nil {
		return "", fmt.Errorf("[window] [GetCookie] [error]: %v", err)
	}
	// wp.Logger.LogString(fmt.Sprintf("[GetCookie] cookies: %s", ck.String()))

	var val string
	cks := strings.Split(tv.String(), ";")
	for _, c := range cks {
		if c != "" {
			// wp.Logger.LogString(fmt.Sprintf("[GetCookie] cookie [%s]", c))
			v := strings.Split(c, "=")
			key := strings.TrimSpace(v[0])
			value := strings.TrimSpace(v[1])
			// wp.Logger.LogString(fmt.Sprintf("[GetCookie] cookie [%s]=%s", key, value))
			if key == _name {
				val = value
				break
			}
		}
	}

	return val, nil
}

// SetCookie ...
func (w *Window) SetCookie(_name, _val string) error {
	if err := ValidJSValue(document, w.document); err != nil {
		return fmt.Errorf("[window] [SetCookie] [error]: %v", err)
	}
	w.document.Set(cookie, fmt.Sprintf("%s=%s", _name, _val))
	// if wp.debug {
	// 	c := wp.document.Get("cookie")
	// 	if c.IsUndefined() || c.IsNull() {
	// 		wp.Logger.LogError(errors.New("goweb_error |DEBUG| 'cookie' js.Value undefined/null"))
	// 		return
	// 	}
	//
	// 	wp.Logger.LogString("goweb |DEBUG| current cookies: " + c.String())
	// }
	return nil
}

// GetElementByTag returns the first element in the document with the given tag (e.g. <div>, <a-entity>, etc.)
func (w *Window) GetElementByTag(_tag string) (*Element, error) {
	err := ValidJSValue(document, w.document)
	if err != nil {
		return nil, fmt.Errorf("[window] [GetElementByTag] [error]: %v", err)
	}
	tv := w.document.Call(function__querySelector, _tag)
	if err = ValidJSValue(_tag, tv); err != nil {
		return nil, fmt.Errorf("[window] [GetElementByTag] [error]: %v", err)
	}
	return NewElement(tv), nil
}

// RemoveElementById ...
func (w *Window) RemoveElementById(_id string) {
	err := ValidJSValue(document, w.document)
	if err != nil {
		w.Logger.Error(fmt.Errorf("[window] [RemoveElementById] [error]: %v", err))
		return
	}
	tv := w.document.Call(function__getElementById, _id)
	if err = ValidJSValue(_id, tv); err != nil {
		w.Logger.Error(fmt.Errorf("[window] [RemoveElementByID] [error]: %v", err))
		return
	}
	tv.Call(function__remove)
}
