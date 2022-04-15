//+build tinygo wasm,js

package aframe

type Component interface {
	Attributes() (map[string]interface{}, error)
}

type AComponent struct {
	Component
}
