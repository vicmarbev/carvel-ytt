// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yamlmeta

import (
	"encoding/json"
	"fmt"

	"github.com/k14s/ytt/pkg/filepos"
	"github.com/k14s/ytt/pkg/orderedmap"
	"github.com/k14s/ytt/pkg/yamlmeta/internal/yaml.v2"
)

func (ds *DocumentSet) GetPosition() *filepos.Position { return ds.Position }
func (d *Document) GetPosition() *filepos.Position     { return d.Position }
func (m *Map) GetPosition() *filepos.Position          { return m.Position }
func (mi *MapItem) GetPosition() *filepos.Position     { return mi.Position }
func (a *Array) GetPosition() *filepos.Position        { return a.Position }
func (ai *ArrayItem) GetPosition() *filepos.Position   { return ai.Position }

func (ds *DocumentSet) SetPosition(position *filepos.Position) { ds.Position = position }
func (d *Document) SetPosition(position *filepos.Position)     { d.Position = position }
func (m *Map) SetPosition(position *filepos.Position)          { m.Position = position }
func (mi *MapItem) SetPosition(position *filepos.Position)     { mi.Position = position }
func (a *Array) SetPosition(position *filepos.Position)        { a.Position = position }
func (ai *ArrayItem) SetPosition(position *filepos.Position)   { ai.Position = position }

// DisplayName is used to return a display name for a DocumentSet
func (ds *DocumentSet) DisplayName() string { return "document set" }

// DisplayName is used to return a display name for a Document
func (d *Document) DisplayName() string { return "document" }

// DisplayName is used to return a display name for a Map
func (m *Map) DisplayName() string { return "map" }

// DisplayName is used to return a display name for a MapItem
func (mi *MapItem) DisplayName() string { return "map item" }

// DisplayName is used to return a display name for an Array
func (a *Array) DisplayName() string { return "array" }

// DisplayName is used to return a display name for an ArrayItem
func (ai *ArrayItem) DisplayName() string { return "array item" }

func (ds *DocumentSet) SetValue(val interface{}) error {
	return fmt.Errorf("cannot set value on a %s", ds.DisplayName())
}

func (d *Document) SetValue(val interface{}) error {
	d.ResetValue()
	return d.AddValue(val)
}

func (m *Map) SetValue(val interface{}) error {
	return fmt.Errorf("cannot set value on a %s", m.DisplayName())
}

func (mi *MapItem) SetValue(val interface{}) error {
	mi.ResetValue()
	return mi.AddValue(val)
}

func (a *Array) SetValue(val interface{}) error {
	return fmt.Errorf("cannot set value on an %s", a.DisplayName())
}

func (ai *ArrayItem) SetValue(val interface{}) error {
	ai.ResetValue()
	return ai.AddValue(val)
}

func isValidValue(val interface{}) bool {
	switch val.(type) {
	case *Map, *orderedmap.Map,
		*Array, []interface{},
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool,
		string,
		nil:
		return true
	default:
		return false
	}
}

func (ds *DocumentSet) ResetValue() { ds.Items = nil }
func (d *Document) ResetValue()     { d.Value = nil }
func (m *Map) ResetValue()          { m.Items = nil }
func (mi *MapItem) ResetValue()     { mi.Value = nil }
func (a *Array) ResetValue()        { a.Items = nil }
func (ai *ArrayItem) ResetValue()   { ai.Value = nil }

func (ds *DocumentSet) AddValue(val interface{}) error {
	if item, ok := val.(*Document); ok {
		ds.Items = append(ds.Items, item)
		return nil
	}
	return fmt.Errorf("cannot add non-document value (%T) into documentset", val)
}

func (d *Document) AddValue(val interface{}) error {
	if !isValidValue(val) {
		return fmt.Errorf("documents can only contain arrays, maps, or scalars; this is a %T", val)
	}
	d.Value = val
	return nil
}

func (m *Map) AddValue(val interface{}) error {
	if item, ok := val.(*MapItem); ok {
		m.Items = append(m.Items, item)
		return nil
	}
	return fmt.Errorf("cannot add non-map-item value (%T) into map", val)
}

func (mi *MapItem) AddValue(val interface{}) error {
	if !isValidValue(val) {
		return fmt.Errorf("map items can only contain arrays, maps, or scalars; this is a %T", val)
	}
	mi.Value = val
	return nil
}

func (a *Array) AddValue(val interface{}) error {
	if item, ok := val.(*ArrayItem); ok {
		a.Items = append(a.Items, item)
		return nil
	}
	return fmt.Errorf("cannot add non-array-item value (%T) into array", val)
}

func (ai *ArrayItem) AddValue(val interface{}) error {
	if !isValidValue(val) {
		return fmt.Errorf("array items can only contain maps, arrays, or scalars; this is a %T", val)
	}
	ai.Value = val
	return nil
}

func (ds *DocumentSet) GetValues() []interface{} {
	var result []interface{}
	for _, item := range ds.Items {
		result = append(result, item)
	}
	return result
}

func (d *Document) GetValues() []interface{} { return []interface{}{d.Value} }

func (m *Map) GetValues() []interface{} {
	var result []interface{}
	for _, item := range m.Items {
		result = append(result, item)
	}
	return result
}

func (mi *MapItem) GetValues() []interface{} { return []interface{}{mi.Value} }

func (a *Array) GetValues() []interface{} {
	var result []interface{}
	for _, item := range a.Items {
		result = append(result, item)
	}
	return result
}

func (ai *ArrayItem) GetValues() []interface{} { return []interface{}{ai.Value} }

func (ds *DocumentSet) GetComments() []*Comment { return ds.Comments }
func (d *Document) GetComments() []*Comment     { return d.Comments }
func (m *Map) GetComments() []*Comment          { return m.Comments }
func (mi *MapItem) GetComments() []*Comment     { return mi.Comments }
func (a *Array) GetComments() []*Comment        { return a.Comments }
func (ai *ArrayItem) GetComments() []*Comment   { return ai.Comments }

func (ds *DocumentSet) addComments(comment *Comment) { ds.Comments = append(ds.Comments, comment) }
func (d *Document) addComments(comment *Comment)     { d.Comments = append(d.Comments, comment) }
func (m *Map) addComments(comment *Comment) {
	panic(fmt.Sprintf("Attempted to attach comment (%s) to Map (%v); maps cannot carry comments", comment.Data, m))
}
func (mi *MapItem) addComments(comment *Comment) { mi.Comments = append(mi.Comments, comment) }
func (a *Array) addComments(comment *Comment) {
	panic(fmt.Sprintf("Attempted to attach comment (%s) to Array (%v); arrays cannot carry comments", comment.Data, a))
}
func (ai *ArrayItem) addComments(comment *Comment) { ai.Comments = append(ai.Comments, comment) }

func (ds *DocumentSet) GetAnnotations() interface{} { return ds.annotations }
func (d *Document) GetAnnotations() interface{}     { return d.annotations }
func (m *Map) GetAnnotations() interface{}          { return m.annotations }
func (mi *MapItem) GetAnnotations() interface{}     { return mi.annotations }
func (a *Array) GetAnnotations() interface{}        { return a.annotations }
func (ai *ArrayItem) GetAnnotations() interface{}   { return ai.annotations }

func (ds *DocumentSet) SetAnnotations(anns interface{}) { ds.annotations = anns }
func (d *Document) SetAnnotations(anns interface{})     { d.annotations = anns }
func (m *Map) SetAnnotations(anns interface{})          { m.annotations = anns }
func (mi *MapItem) SetAnnotations(anns interface{})     { mi.annotations = anns }
func (a *Array) SetAnnotations(anns interface{})        { a.annotations = anns }
func (ai *ArrayItem) SetAnnotations(anns interface{})   { ai.annotations = anns }

// Below methods disallow marshaling of nodes directly
var _ []yaml.Marshaler = []yaml.Marshaler{&DocumentSet{}, &Document{}, &Map{}, &MapItem{}, &Array{}, &ArrayItem{}}

func (ds *DocumentSet) MarshalYAML() (interface{}, error) { panic("Unexpected marshaling of docset") }
func (d *Document) MarshalYAML() (interface{}, error)     { panic("Unexpected marshaling of doc") }
func (m *Map) MarshalYAML() (interface{}, error)          { panic("Unexpected marshaling of map") }
func (mi *MapItem) MarshalYAML() (interface{}, error)     { panic("Unexpected marshaling of mapitem") }
func (a *Array) MarshalYAML() (interface{}, error)        { panic("Unexpected marshaling of array") }
func (ai *ArrayItem) MarshalYAML() (interface{}, error)   { panic("Unexpected marshaling of arrayitem") }

// Below methods disallow marshaling of nodes directly
var _ []json.Marshaler = []json.Marshaler{&DocumentSet{}, &Document{}, &Map{}, &MapItem{}, &Array{}, &ArrayItem{}}

func (ds *DocumentSet) MarshalJSON() ([]byte, error) { panic("Unexpected marshaling of docset") }
func (d *Document) MarshalJSON() ([]byte, error)     { panic("Unexpected marshaling of doc") }
func (m *Map) MarshalJSON() ([]byte, error)          { panic("Unexpected marshaling of map") }
func (mi *MapItem) MarshalJSON() ([]byte, error)     { panic("Unexpected marshaling of mapitem") }
func (a *Array) MarshalJSON() ([]byte, error)        { panic("Unexpected marshaling of array") }
func (ai *ArrayItem) MarshalJSON() ([]byte, error)   { panic("Unexpected marshaling of arrayitem") }

func (ds *DocumentSet) sealed() {}
func (d *Document) sealed()     {}
func (m *Map) sealed()          {}
func (mi *MapItem) sealed()     {}
func (a *Array) sealed()        {}
func (ai *ArrayItem) sealed()   {}

func (n *DocumentSet) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *DocumentSet) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
func (n *Document) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *Document) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
func (n *Map) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *Map) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
func (n *MapItem) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *MapItem) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
func (n *Array) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *Array) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
func (n *ArrayItem) GetMeta(name string) interface{} {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	return n.meta[name]
}
func (n *ArrayItem) SetMeta(name string, data interface{}) {
	if n.meta == nil {
		n.meta = make(map[string]interface{})
	}
	n.meta[name] = data
}
