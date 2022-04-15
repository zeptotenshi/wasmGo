//+build tinygo wasm,js

package firebase

import (
	"fmt"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	firebase = "firebase"

	ERROR__code    = "code"
	ERROR__message = "message"
)

type Firebase struct {
	value js.Value
}

func NewClient() (*Firebase, error) {
	fb := js.Global().Get(firebase)
	if err := web.ValidJSValue(firebase, fb); err != nil {
		return nil, err
	}
	return &Firebase{value: fb}, nil
}

func (f *Firebase) Auth() (*Auth, error) {
	if err := web.ValidJSValue(firebase, f.value); err != nil {
		return nil, err
	}
	authClient := f.value.Call(auth)
	if err := web.ValidJSValue(fmt.Sprintf("%s.%s", firebase, auth), authClient); err != nil {
		return nil, err
	}
	return &Auth{value: authClient, User: js.ValueOf(nil)}, nil
}

func (f *Firebase) Store() (*Firestore, error) {
	if err := web.ValidJSValue(firebase, f.value); err != nil {
		return nil, err
	}
	storeClient := f.value.Call(firestore)
	if err := web.ValidJSValue(fmt.Sprintf("%s.%s", firebase, firestore), storeClient); err != nil {
		return nil, err
	}
	return &Firestore{value: storeClient}, nil
}

func (f *Firebase) Storage() (*Storage, error) {
	if err := web.ValidJSValue(firebase, f.value); err != nil {
		return nil, err
	}
	storageClient := f.value.Call(storage)
	if err := web.ValidJSValue(fmt.Sprintf("%s.%s", firebase, storage), storageClient); err != nil {
		return nil, err
	}
	return &Storage{value: storageClient}, nil
}
