package goluaez

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
)

type State struct {
	mu sync.Mutex
	*lua.LState
}

// dont modify this directly
// modify code in data/prelude.lua and then put it here.
const prelude = `
-- remove whitespace at the ends of a string.
function string:strip()
    return self:match'^%s*(.*%S)' or ''
end

-- split a string by a separator
function string:split(sep)
    local sep, fields = sep or "\t", {}
    local pattern = string.format("[^%s]+", sep)
    for tok in self:gmatch(pattern) do fields[#fields+1] = tok end
    return fields
end
`

// TODO: see https://github.com/yuin/gluamapper

// NewState creates a new State object optionally initialized with some code.
func NewState(code ...string) (*State, error) {
	s := &State{}
	options := lua.Options{IncludeGoStackTrace: true}
	s.LState = lua.NewState(options)
	s.PreloadModule("re", gluare.Loader)
	var err error
	err = s.LState.DoString(prelude)
	if err != nil {
		return nil, err
	}
	if len(code) != 0 && len(code[0]) != 0 {
		err = s.LState.DoString(code[0])
		if err != nil {
			return s, err
		}
	}
	return s, err
}

func (s *State) SetGlobal(name string, val interface{}) error {
	l, err := Go2LValue(val)
	s.LState.SetGlobal(name, l)
	return err
}

func Go2LValue(v interface{}) (lua.LValue, error) {
	switch cast := v.(type) {
	case float32:
		return lua.LNumber(cast), nil
	case float64:
		return lua.LNumber(cast), nil
	case int, int32, int64:
		return lua.LNumber(float64(reflect.ValueOf(v).Int())), nil
	case uint, uint32, uint64:
		return lua.LNumber(float64(reflect.ValueOf(v).Uint())), nil
	case string:
		return lua.LString(cast), nil
	case bool:
		return lua.LBool(cast), nil
	case nil:
		return lua.LNil, nil
	case []string:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LString(val))
		}
		return tbl, nil
	case []int:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LNumber(val))
		}
		return tbl, nil
	case []int64:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LNumber(val))
		}
		return tbl, nil
	case []int32:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LNumber(val))
		}
		return tbl, nil
	case []float32:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LNumber(val))
		}
		return tbl, nil
	case []float64:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LNumber(val))
		}
		return tbl, nil
	case []bool:
		tbl := &lua.LTable{}
		for _, val := range cast {
			tbl.Append(lua.LBool(val))
		}
		return tbl, nil
	default:
		return nil, fmt.Errorf("cant convert %v to LValue", v)
	}
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
}

// Run code given some values. This is thread-safe.
func (s *State) Run(code string, values ...map[string]interface{}) (interface{}, error) {
	var err error
	if len(values) != 0 {
		lvals := make([]lua.LValue, len(values[0]))
		keys := make([]string, len(values[0]))
		j := 0
		for k, v := range values[0] {
			s.mu.Lock()
			lvals[j], err = Go2LValue(v)
			s.mu.Unlock()
			if err != nil {
				return nil, err
			}
			keys[j] = k
			j++
		}

		s.mu.Lock()
		for i, lv := range lvals {
			s.LState.SetGlobal(keys[i], lv)
		}
	} else {
		s.mu.Lock()
	}
	if err := s.DoString("return " + code); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	if s.GetTop() == 0 {
		s.mu.Unlock()
		return nil, nil
	}
	v := s.Get(-1)
	s.Pop(1)
	s.mu.Unlock()
	return LValue2Go(v)
}
