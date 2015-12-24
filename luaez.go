package goluaez

import (
	"log"
	"reflect"
	"sync"

	"github.com/layeh/gopher-luar"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
)

type State struct {
	mu sync.Mutex
	*lua.LState
}

// TODO: see https://github.com/yuin/gluamapper

// NewState creates a new State object optionally initialized with some code.
func NewState(code ...string) (*State, error) {
	s := &State{}
	s.LState = lua.NewState()
	s.PreloadModule("re", gluare.Loader)
	var err error
	if len(code) != 0 && len(code[0]) != 0 {
		err = s.LState.DoString(code[0])
		if err != nil {
			return s, err
		}
	}
	return s, err
}

func (s *State) SetGlobal(name string, val interface{}) {
	s.LState.SetGlobal(name, luar.New(s.LState, val))
}

func LValue2Go(v lua.LValue) (interface{}, error) {
	switch v.Type() {
	case lua.LTString:
		return string(v.(lua.LString)), nil
	case lua.LTNumber:
		return float64(v.(lua.LNumber)), nil
	case lua.LTBool:
		return bool(v.(lua.LBool)), nil
	case lua.LTNil:
		return nil, nil
	case lua.LTTable:
		tbl := v.(*lua.LTable)
		varr := make([]interface{}, 0)
		vmap := make(map[string]interface{})
		all_ints := true
		k := lua.LNil
		i := 0
		for {
			i += 1
			key, val := tbl.Next(k)
			if key == lua.LNil {
				break
			}
			gokey, err := LValue2Go(key)
			if err != nil {
				return nil, err
			}
			goval, err := LValue2Go(val)
			if err != nil {
				return nil, err
			}

			// see if we have all int keys
			if gofloat, ok := gokey.(float64); ok && all_ints {
				if int(gofloat) == i {
					varr = append(varr, goval)
				} else {
					all_ints = false
				}
			} else {
				all_ints = false
			}
			vmap[key.String()] = goval
			k = key
		}
		if all_ints {
			return varr, nil
		}
		return vmap, nil
	case lua.LTUserData:
		goval := v.(*lua.LUserData).Value
		return goval, nil
	default:
		switch t := v.(type) {
		default:
			log.Println("IN luaez ...", t)
			log.Printf("type:%+v\n", v)
			log.Println(v.(*lua.LUserData).Value)
			return reflect.ValueOf(t), nil
		}

	}
	return v, nil

}

func (s *State) Run(code string, values ...map[string]interface{}) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(values) != 0 {
		for k, v := range values[0] {
			s.SetGlobal(k, luar.New(s.LState, v))
		}
	}
	if err := s.DoString("return " + code); err != nil {
		return nil, err
	}
	if s.GetTop() == 0 {
		return nil, nil
	}
	v := s.Get(-1)
	s.Pop(1)
	return LValue2Go(v)
}
