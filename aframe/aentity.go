//+build tinygo wasm,js

package aframe

import (
	"fmt"
	"math"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	entity__tag = "a-entity"

	PROPERTY__components = "components"
	PROPERTY__data       = "data"
	PROPERTY__isEntity   = "isEntity"

	PROPERTY__children = "children"
	PROPERTY__parentEl = "parentEl"

	PROPERTY__geometry = "geometry"
	PROPERTY__length   = "length"
	PROPERTY__width    = "width"
	PROPERTY__depth    = "depth"

	PROPERTY__material = "material"
	PROPERTY__color    = "color"
	PROPERTY__emissive = "emissive"
	PROPERTY__src      = "src"

	PROPERTY__text  = "text"
	PROPERTY__value = "value"

	PROPERTY__point = "point"

	function__destroy = "destroy"
)

type AEntity struct {
	*web.Element
	scene *Aframe
}

func (e *AEntity) Scene() *Aframe {
	return e.scene
}

func (e *AEntity) AppendChild(_ae *AEntity) error {
	return e.Element.SetChild(_ae.Element.Value)
}

func (e *AEntity) SetPosition(_x, _y, _z float64) error {
	position, err := e.Element.GetProperty(PROPERTY__object3D, PROPERTY__position)
	if err != nil {
		return fmt.Errorf("[AEntity] %s [SetPosition] [error]: %v", e.Element, err)
	}
	position.Set(PROPERTY__x, _x)
	position.Set(PROPERTY__y, _y)
	position.Set(PROPERTY__z, _z)
	return nil
}

func (e *AEntity) SetRotation(_x, _y, _z float64) error {
	rotation, err := e.Element.GetProperty(PROPERTY__object3D, PROPERTY__rotation)
	if err != nil {
		return fmt.Errorf("[AEntity] %s [SetRotation] [error]: %v", e.Element, err)
	}
	rotation.Set(PROPERTY__x, _x*math.Pi/180)
	rotation.Set(PROPERTY__y, _y*math.Pi/180)
	rotation.Set(PROPERTY__z, _z*math.Pi/180)
	return nil
}

func (e *AEntity) SetVisible(_on bool) error {
	obj, err := e.Element.GetProperty(PROPERTY__object3D)
	if err != nil {
		return fmt.Errorf("[AEntity] %s [setVisible] [error]: %v", e.Element, err)
	}
	obj.Set(PROPERTY__visible, _on)
	return nil
}

func (e *AEntity) Append() error {
	if err := web.ValidJSValue(e.Element.String(), e.Element.Value); err != nil {
		return fmt.Errorf("[AEntity] %s [append] [error]: %v", e.Element, err)
	}
	if e.scene == nil {
		return fmt.Errorf("[AEntity] %s [append] [error]: scene ref nil", e.Element)
	}
	if err := web.ValidJSValue("scene", e.scene.scene); err != nil {
		return fmt.Errorf("[AEntity] %s [append] [error]: %v", e.Element, err)
	}

	e.scene.scene.Call(web.FUNCTION__appendChild, e.Element.Value)

	return nil
}

func (e *AEntity) Remove(_wc bool) {
	isEnt, err := e.Element.GetProperty(PROPERTY__isEntity)
	if err == nil {
		if isEnt.Bool() {
			if _wc {
				e.RemoveChildren()
			}

			if err := e.Element.Remove(); err != nil {
				e.scene.Error(fmt.Errorf("[AEntity] %s [Remove] [error]: %v", e.Element, err))
				return
			}
			e.Element.Value.Call(function__destroy)
		}

	} else {
		if e.Element.Tag == entity__tag {
			e.scene.Error(fmt.Errorf("[AEntity] %s [Remove] [error]: has tag <a-entity> but isEnt[false]", e.Element))
		}
		e.Element.Remove()
	}

	delete(e.scene.entities, e.Element.ID)
}

func (e *AEntity) RemoveChildren() {
	cl, err := e.Element.GetProperty(PROPERTY__children)
	if err != nil {
		e.scene.Error(fmt.Errorf("[AEntity] %s [RemoveChildren] [error]: %v", e.Element, err))
		return
	}
	for i := 0; i < cl.Length(); i++ {
		c := cl.Index(i)
		cid := c.Get(web.ELEMENT__id).String()
		if cid != "" && e.scene != nil {
			ce := e.scene.GetEntityByID(cid)
			if ce == nil {
				e.scene.Error(fmt.Errorf("[AEntity] %s [RemoveChildren] [error]: child[%s] nil", e.Element, cid))
				continue
			}
			ce.Remove(true)
		}
	}
}
