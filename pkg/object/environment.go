package object

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: outer}
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) GetAllBuiltinInfos() map[string]BuiltIn {
	// Reduce the need to convert the type again
	infos := map[string]BuiltIn{}
	// Make sure the keys reponse in specific order
	keys := []string{}
	for k, v := range e.store {
		if bi, ok := v.(BuiltIn); ok {
			keys = append(keys, k)
			infos[k] = bi
		}
	}
	return infos
}

func (e *Environment) GetAllStoreData() string {
	keys := []string{}
	for k := range e.store {
		keys = append(keys, k)
	}
	// Make sure the order is good enough
	sort.Strings(keys)
	slices.SortFunc(keys, func(a, b string) int {
		return strings.Compare(string(e.store[a].Type()), string(e.store[b].Type()))
	})

	var helpInfo strings.Builder

	fmt.Fprintf(&helpInfo, "Store data:")
	for i, k := range keys {
		if i > 0 {
			fmt.Fprintf(&helpInfo, "\n")
		}
		fmt.Fprintf(&helpInfo, "\t- %s (%s): %s", k, e.store[k].Type(), e.store[k].Inspect())
	}
	return helpInfo.String()
}
