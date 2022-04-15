//+build tinygo wasm,js

package aframe

import (
	"fmt"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	texture  = "texture"
	material = "material"

	property__side = "side"

	function__load = "load"

	Skybox__front  = "front"
	Skybox__back   = "back"
	Skybox__left   = "left"
	Skybox__right  = "right"
	Skybox__top    = "top"
	Skybox__bottom = "bottom"
)

type Skybox struct {
	Name string
	uuid int

	Images map[string]string

	Length float32
	Height float32
	Depth  float32
}

func (af *Aframe) newTexture(_src string) (js.Value, error) {
	tv := js.ValueOf(nil)
	err := web.ValidJSValue(THREE__TextureLoader, af.Three.TextureLoader)
	if err != nil {
		return tv, fmt.Errorf("[aframe] [newtexture] [error]: %v", err)
	}
	textureLoader := af.Three.TextureLoader.New()
	tv = textureLoader.Call(function__load, _src)
	if err = web.ValidJSValue(texture, tv); err != nil {
		return tv, fmt.Errorf("[aframe] [newtexture] [error]: %v", err)
	}
	return tv, nil
}

func (af *Aframe) newBasicMaterial(_texture js.Value) (js.Value, error) {
	tv := js.ValueOf(nil)
	err := web.ValidJSValue(THREE__MeshBasicMaterial, af.Three.MeshBasicMaterial)
	if err != nil {
		return tv, fmt.Errorf("[aframe] [newbasicmaterial] [error]: %v", err)
	}
	tv = af.Three.MeshBasicMaterial.New(map[string]interface{}{"map": _texture})
	if err = web.ValidJSValue(material, tv); err != nil {
		return tv, fmt.Errorf("[aframe] [newbasicmaterial] [error]: %v", err)
	}
	return tv, nil
}

func (af *Aframe) SetSkybox(_sky *Skybox) error {
	err := web.ValidJSValue(scene, af.scene)
	if err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	if err = web.ValidJSValue(THREE, af.Three.Value); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	if err = web.ValidJSValue(THREE__BackSide, af.Three.BackSide); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	if err = web.ValidJSValue(THREE__Mesh, af.Three.Mesh); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	if err = web.ValidJSValue(THREE__BoxGeometry, af.Three.BoxGeometry); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}

	if l := len(_sky.Images); l != 6 {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: Skybox[%s] image map does not have 6 images", _sky.Name)
	}

	materialArray := make([]interface{}, 6)

	for i, fn := range []string{Skybox__front, Skybox__back, Skybox__top, Skybox__bottom, Skybox__right, Skybox__left} {
		v, ok := _sky.Images[fn]
		if !ok {
			return fmt.Errorf("[aframe] [SetSkybox] [error]: face[%s] image not found in image map", fn)
		}
		texture, err := af.newTexture(v)
		if err != nil {
			return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
		}
		material, err := af.newBasicMaterial(texture)
		if err != nil {
			return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
		}
		material.Set(property__side, af.Three.BackSide)

		materialArray[i] = material
	}

	skyboxGeo := af.Three.BoxGeometry.New(_sky.Length, _sky.Height, _sky.Depth)
	if err = web.ValidJSValue(THREE__BoxGeometry, skyboxGeo); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	skyboxMesh := af.Three.Mesh.New(skyboxGeo, materialArray)
	if err = web.ValidJSValue(THREE__Mesh, skyboxMesh); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	id := skyboxMesh.Get(web.ELEMENT__id)
	if err = web.ValidJSValue(web.ELEMENT__id, id); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}

	sceneObj := af.scene.Get(PROPERTY__object3D)
	if err = web.ValidJSValue(fmt.Sprintf("%s.%s", scene, PROPERTY__object3D), sceneObj); err != nil {
		return fmt.Errorf("[aframe] [SetSkybox] [error]: %v", err)
	}
	sceneObj.Call(function__add, skyboxMesh)

	_sky.uuid = id.Int()
	af.skyboxes[_sky.Name] = id.Int()

	return nil
}
