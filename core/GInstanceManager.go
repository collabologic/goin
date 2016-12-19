// Package: Framework package in GoIn.
package core

import (
	"reflect"
)

// Factory function type for GInstanceManager.
type FactoryFunc func() interface{}

// Instance of GInstanceManager.
var gInstanceManager *GInstanceManager 
// Flag of created GInstanceManager.
var created bool

// Factory function to GInstanceManager
func GetGInstanceManager() GInstanceManager{
	if !created {
		gInstanceManager = new(GInstanceManager)
		gInstanceManager.factories = make(map[string]FactoryFunc)
		gInstanceManager.singletons = make(map[string]interface{})
		created = true
	}
	return *gInstanceManager
}

// Struct of InstanceManager
type GInstanceManager struct {
	factories map[string]FactoryFunc
	singletons map[string]interface{}
}

// Get or create instance.
// If struct is used singleton, then this function call.
func (gic *GInstanceManager) Get(strName string) (interface{}) {
	if _, ok := gic.singletons[strName]; !ok {
		newSingleton := gic.factories[strName]()
		autoInjection(gic, newSingleton)
		gic.singletons[strName] = newSingleton
	}
	return gic.singletons[strName]
}

// Create instance.
func (gic *GInstanceManager) New(strName string) (interface{}) {
	newObj := gic.factories[strName]()
	newObj = autoInjection(gic, newObj)
	return newObj
}

// Regist factory method to GInstanceManager.
func (gic *GInstanceManager) AddFactoryMethod(strName string, fn FactoryFunc) {
	if _, ok := gic.factories[strName] ; ok{
		panic("GIManager.Duplicate.Factory")
	}
	gic.factories[strName] = fn
}

// Process of struct instance creating.
// Resolve GITag.
func autoInjection( gic *GInstanceManager, target interface{} ) interface{} {
	value := reflect.ValueOf(target)
	elm := value.Elem() // 渡されているインターフェイスにはポインタが入っているので実体を得る
	typ := elm.Type()
	if typ.Kind() != reflect.Struct {
		return target
	}
	for f := 0 ;  f<=typ.NumField() ; f++ {
		field := typ.Field(f)
		var val interface{}
		// wiredタグの処理
		if wired := field.Tag.Get("wired") ; wired != "" {
			if _, ok := gic.singletons[wired]; !ok {
				val = gic.factories[wired]()
				val = autoInjection(gic, val)
			} else {
				val = gic.singletons[wired] 
			}
		}
		// injectタグの処理
		if inject := field.Tag.Get("inject") ; inject != "" {
			val = gic.factories[inject]()
			val = autoInjection(gic, val)
		}
		if val == nil {
			continue
		}
		elm.Field(f).Set(reflect.ValueOf(val))
	}
	return value.Interface()
}
