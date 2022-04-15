//+build tinygo wasm,js

package firebase

import (
	"fmt"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	firestore = "firestore"

	function__doc = "doc"
)

type Firestore struct {
	value js.Value
}

func (f *Firestore) GetDoc(_p string) (js.Value, error) {
	err := web.ValidJSValue(firestore, f.value)
	if err != nil {
		return js.ValueOf(nil), fmt.Errorf("[firebase] [firestore] [GetDoc] [error]: %v", err)
	}

	doc := f.value.Call(function__doc, _p)
	if err = web.ValidJSValue(fmt.Sprintf("%s.%s", function__doc, _p), doc); err != nil {
		return js.ValueOf(nil), fmt.Errorf("[firebase] [firestore] [GetDoc] [error]: %v", err)
	}

	return doc, nil
}
