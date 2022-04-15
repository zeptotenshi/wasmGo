//+build tinygo wasm,js

package aframe

import (
	"syscall/js"
)

const (
	THREE = "THREE"

	THREE__Vector3              = "Vector3"
	THREE__BackSide             = "BackSide"
	THREE__BoxGeometry          = "BoxGeometry"
	THREE__ConeGeometry         = "ConeGeometry"
	THREE__CircleGeometry       = "CircleGeometry"
	THREE__RingGeometry         = "RingGeometry"
	THREE__TextureLoader        = "TextureLoader"
	THREE__Mesh                 = "Mesh"
	THREE__MeshBasicMaterial    = "MeshBasicMaterial"
	THREE__MeshStandardMAterial = "MeshStandardMaterial"

	PROPERTY__object3D = "object3D"
	PROPERTY__position = "position"
	PROPERTY__rotation = "rotation"
	PROPERTY__visible  = "visible"
	PROPERTY__x        = "x"
	PROPERTY__y        = "y"
	PROPERTY__z        = "z"
)

type Three struct {
	js.Value
	Vector3              js.Value
	BackSide             js.Value
	BoxGeometry          js.Value
	TextureLoader        js.Value
	Mesh                 js.Value
	MeshBasicMaterial    js.Value
	MeshStandardMaterial js.Value

	// ConeGeometry         js.Value
	// CircleGeometry       js.Value
	// RingGeometry         js.Value
}

func NewThree(_v js.Value) *Three {
	return &Three{
		Value:                _v,
		Vector3:              _v.Get(THREE__Vector3),
		BackSide:             _v.Get(THREE__BackSide),
		BoxGeometry:          _v.Get(THREE__BoxGeometry),
		TextureLoader:        _v.Get(THREE__TextureLoader),
		Mesh:                 _v.Get(THREE__Mesh),
		MeshBasicMaterial:    _v.Get(THREE__MeshBasicMaterial),
		MeshStandardMaterial: _v.Get(THREE__MeshStandardMAterial),
		// ConeGeometry:         _v.Get("ConeGeometry"),
		// CircleGeometry:       _v.Get("CircleGeometry"),
		// RingGeometry:         _v.Get("RingGeometry"),
	}
}
