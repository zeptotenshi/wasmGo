//+build tinygo wasm,js

package firebase

import (
	"syscall/js"
)

const (
	storage = "storage"
)

type Storage struct {
	value js.Value
}
