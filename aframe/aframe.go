//+build tinygo wasm,js

package aframe

import (
	// "fmt"
	"strings"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	aframe = "AFRAME"
	scene  = "scene"

	element__scene = "a-scene"

	function__add = "add"
)

type Aframe struct {
	js.Value
	scene js.Value
	*web.Window
	Three *Three

	entities map[string]*AEntity
	skyboxes map[string]int
}

func NewAframe(_wp *web.Window) *Aframe {
	af := &Aframe{
		Window:   _wp,
		entities: map[string]*AEntity{},
		skyboxes: map[string]int{},
	}

	af.Value, _ = _wp.GetGlobal(aframe)

	scene, _ := _wp.GetElementByTag(element__scene)
	af.scene = scene.Value

	three, _ := _wp.GetGlobal(THREE)
	af.Three = NewThree(three)

	// v.Call("registerComponent", "genesis", map[string]interface{}{
	// 	"schema": map[string]interface{}{"guardian": map[string]interface{}{"default": ""}},
	// 	"init": js.FuncOf(func(_this js.Value, _args []js.Value) interface{} {
	// 		_this.Set("initialized", true)

	// 		// var tempString string
	// 		// var tempValue js.Valu

	// 		// tempString = _this.Get("data").Get("guardian").String()

	// 		// if tempString != "" {

	// 		// }

	// 		return js.ValueOf(nil)
	// 	}),
	// })

	return af
}

func (af *Aframe) NewEntity(_el *web.Element) *AEntity {
	r := &AEntity{
		Element: _el,
		scene:   af,
	}
	af.entities[_el.ID] = r
	return r
}

func (af *Aframe) NewEntityWithID(_id string) *AEntity {
	tempEl := af.Window.NewElementWithTag(entity__tag)
	tempEl.SetID(_id)

	r := af.NewEntity(tempEl)
	// if err := r.SetAttribute("genesis", map[string]interface{}{}); err != nil {
	// 	af.Error(fmt.Errorf("[AFrame|NewEntityWithID] error: %v", err))
	// }

	// af.Info(fmt.Sprintf("[Aframe|NewEntityWithID] make Entity[%s]", _id))

	return r
}

func (af *Aframe) GetEntityByID(_id string) *AEntity {
	ent, ok := af.entities[_id]
	if !ok {
		el := af.Window.ElementById(_id)
		if el == nil {
			ent = af.NewEntityWithID(_id)
		} else {
			ent = af.NewEntity(el)
		}
		af.entities[_id] = ent
	}
	return ent
}

func (af *Aframe) RemoveEntityByID(_id string) {
	tempEntity, entityExists := af.entities[_id]
	if !entityExists {
		return
	}
	tempEntity.Remove(true)
	delete(af.entities, _id)
}

func (af *Aframe) GetEntitiesByClass(_className string) []*AEntity {
	es := []*AEntity{}
	// af.Debug(fmt.Sprintf("[Aframe|GetEntitiesByClass] entity count[%d]", len(af.entities)))
	for _, e := range af.entities {
		cn := e.Value.Get(web.ELEMENT__class).String()
		if cn == "" {
			continue
		}
		// af.Debug(fmt.Sprintf("[Aframe|GetEntitiesByClass] entity class[%s]", cn))
		if strings.Contains(cn, _className) {
			es = append(es, e)
		}
	}

	return es
}
