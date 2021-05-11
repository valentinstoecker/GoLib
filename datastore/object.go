package datastore

import (
	"encoding/gob"
	"fmt"
	"os"
)

// ObjectID (id of object)
type ObjectID uint64

// AttrID (id of attribute)
type AttrID uint64

// RelID (id of relation)
type RelID uint64

type ds struct {
	CObjID   ObjectID
	CAttrID  AttrID
	CRelID   RelID
	Attr2str map[AttrID]string
	Str2attr map[string]AttrID
	Rel2str  map[RelID]string
	Str2rel  map[string]RelID
	Objects  map[ObjectID]*Object
}

type Request struct {
	DataStore
	pred func(o *Object) bool
}

func NewReq(d DataStore) Request {
	return Request{
		DataStore: d,
		pred:      func(o *Object) bool { return true },
	}
}

func (req Request) StringEqual(name string, value string) Request {
	return Request{
		DataStore: req.DataStore,
		pred: func(o *Object) bool {
			str, ok := req.GetAttribute(o, name).(string)
			return req.pred(o) && ok && str == value
		},
	}
}

type DataStore interface {
	AddObject() *Object
	GetObjects() []*Object
	SetAttribute(obj *Object, name string, data interface{})
	GetAttribute(obj *Object, name string) interface{}
	AddRel(src *Object, name string, tgt *Object)
	GetRel(obj *Object, name string) []*Object
	Print(obj *Object)
	Save(p string) error
	StringEqual(obj *Object, name string, value string) bool
}

func NewDS() DataStore {
	return &ds{
		Attr2str: make(map[AttrID]string),
		Str2attr: make(map[string]AttrID),
		Rel2str:  make(map[RelID]string),
		Str2rel:  make(map[string]RelID),
		Objects:  make(map[ObjectID]*Object),
	}
}

func (d *ds) AddObject() *Object {
	obj := new(Object)
	obj.ID = d.CObjID
	d.Objects[d.CObjID] = obj
	d.CObjID++
	return obj
}

func (d *ds) GetObjects() []*Object {
	objs := make([]*Object, len(d.Objects))
	i := 0
	for _, v := range d.Objects {
		objs[i] = v
		i++
	}
	return objs
}

func (d *ds) SetAttribute(obj *Object, name string, data interface{}) {
	attrID, ok := d.Str2attr[name]
	if !ok {
		attrID = d.CAttrID
		d.CAttrID++
		d.Attr2str[attrID] = name
		d.Str2attr[name] = attrID
	}
	if obj.Attr == nil {
		obj.Attr = make(map[AttrID]interface{})
	}
	obj.Attr[attrID] = data
}

func (d *ds) GetAttribute(obj *Object, name string) interface{} {
	attrID, ok := d.Str2attr[name]
	if !ok {
		return nil
	}
	if obj.Attr == nil {
		return nil
	}
	v, ok := obj.Attr[attrID]
	if !ok {
		return nil
	}
	return v
}

func (d *ds) AddRel(src *Object, name string, tgt *Object) {
	relID, ok := d.Str2rel[name]
	if !ok {
		relID = d.CRelID
		d.CRelID++
		d.Rel2str[relID] = name
		d.Str2rel[name] = relID
	}
	if src.Rels == nil {
		src.Rels = make(map[RelID][]ObjectID)
	}
	src.Rels[relID] = append(src.Rels[relID], tgt.ID)
}

func (d *ds) GetRel(obj *Object, name string) []*Object {
	if id, ok := d.Str2rel[name]; ok {
		if ts, ok := obj.Rels[id]; ok {
			objs := make([]*Object, len(ts))
			for i, v := range ts {
				objs[i] = d.Objects[v]
			}
			return objs
		}
	}
	return []*Object{}
}

func (d *ds) Print(obj *Object) {
	fmt.Println("-= Object =-")
	for id, v := range obj.Attr {
		attrStr := d.Attr2str[id]
		fmt.Println(attrStr+":", v)
	}
	for id, rel := range obj.Rels {
		relStr := d.Rel2str[id]
		fmt.Println(relStr+":", len(rel))
	}
	fmt.Println("-==========-")
}

func (d *ds) Save(p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(d); err != nil {
		return err
	}
	return nil
}

func Open(p string) (DataStore, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	dataStore := NewDS()
	if err := dec.Decode(dataStore); err != nil {
		return nil, err
	}
	return dataStore, nil
}

func GetObjectsWhere(req Request) []*Object {
	res := make([]*Object, 0)
	for _, obj := range req.DataStore.GetObjects() {
		if req.pred(obj) {
			res = append(res, obj)
		}
	}
	return res
}

func (d *ds) StringEqual(obj *Object, name string, value string) bool {
	str, ok := d.GetAttribute(obj, name).(string)
	return ok && str == value
}

type Object struct {
	ID   ObjectID
	Attr map[AttrID]interface{}
	Rels map[RelID][]ObjectID
}
