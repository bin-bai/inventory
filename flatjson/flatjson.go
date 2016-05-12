// flatjson
package flatjson

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	"github.com/bin-bai/inventory"
)

type FlatJson struct {
	Cols map[string]*FlatJsonCol `json:"Cols"`

	path string
	lock sync.RWMutex
}

type FlatJsonCol struct {
	Docs map[string]string `json:"Docs"`

	elemType reflect.Type
	lock     sync.RWMutex
}

func NewFlatJson(path string) *FlatJson {
	fj := new(FlatJson)
	fj.Cols = make(map[string]*FlatJsonCol)
	fj.path = path

	return fj
}

func (fj *FlatJson) Load() error {
	fj.lock.Lock()
	defer fj.lock.Unlock()

	buf, err := ioutil.ReadFile(fj.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &fj.Cols)
	if err != nil {
		return err
	}

	return nil
}

func (fj *FlatJson) Save() error {
	fj.lock.Lock()
	defer fj.lock.Unlock()

	f, err := os.OpenFile(fj.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := json.MarshalIndent(fj.Cols, "", "\t")
	if err != nil {
		return err
	}

	f.Write(buf)
	return nil
}

func (fj *FlatJson) Close() {
	fj.lock.Lock()
	defer fj.lock.Unlock()

	for k, _ := range fj.Cols {
		delete(fj.Cols, k)
	}

	fj.path = ""
}

func (fj *FlatJson) Use(name string, elemtype interface{}) inventory.Collection {
	fj.lock.Lock()
	defer fj.lock.Unlock()

	c, ok := fj.Cols[name]
	if !ok {
		c = new(FlatJsonCol)
		c.Docs = make(map[string]string)
		fj.Cols[name] = c
	}
	c.elemType = reflect.TypeOf(elemtype)

	return c
}

func (fj *FlatJson) ColSize() int {
	return len(fj.Cols)
}

func (jc *FlatJsonCol) Set(key string, elem interface{}) error {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	buf, err := json.Marshal(elem)
	if err != nil {
		return err
	}

	jc.Docs[key] = string(buf)
	return nil
}

func (jc *FlatJsonCol) Get(key string) interface{} {
	if jc.elemType == nil {
		return nil
	}

	jc.lock.RLock()
	defer jc.lock.RUnlock()

	data, ok := jc.Docs[key]
	if !ok {
		return nil
	}

	elem := reflect.New(jc.elemType)
	err := json.Unmarshal([]byte(data), elem.Interface())
	if err != nil {
		return nil
	}

	if !reflect.Indirect(elem).CanInterface() {
		return nil
	}
	return reflect.Indirect(elem).Interface()
}

func (jc *FlatJsonCol) Del(key string) {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	delete(jc.Docs, key)
}

func (jc *FlatJsonCol) GetIndex(i int) (string, interface{}) {
	jc.lock.RLock()
	defer jc.lock.RUnlock()

	for k, _ := range jc.Docs {
		if i == 0 {
			return k, jc.Get(k)
		}
		i--
	}

	return "", nil
}

func (jc *FlatJsonCol) Size() int {
	return len(jc.Docs)
}
