package goluaez_test

import (
	"log"
	"testing"

	"github.com/brentp/goluaez"
	"github.com/yuin/gopher-lua"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Tester struct{}

var _ = Suite(&Tester{})

func (t *Tester) TestNew(c *C) {
	s, err := goluaez.NewState()
	c.Assert(err, IsNil)
	c.Assert(s, Not(IsNil))
}

func (t *Tester) TestNewWithCode(c *C) {
	s, err := goluaez.NewState("fnction(a, b) return a + b end")
	c.Assert(err, Not(IsNil))

	s, err = goluaez.NewState("uu = function(a, b) return a + b end")
	c.Assert(err, IsNil)
	c.Assert(s, Not(IsNil))
}

func (t *Tester) TestFuncFloat(c *C) {
	s, err := goluaez.NewState("uu = function(a, b) return a + b end")
	c.Assert(err, IsNil)
	v, e := s.Run("uu(2, 9)")
	c.Assert(e, IsNil)
	c.Assert(v, Equals, float64(11))
}

func (t *Tester) TestFuncString(c *C) {
	s, err := goluaez.NewState("uu = function(a, b) return a .. b end")
	c.Assert(err, IsNil)
	v, e := s.Run("uu(2, 9)")
	c.Assert(e, IsNil)
	c.Assert(v, Equals, "29")

}

func (t *Tester) TestFuncBool(c *C) {
	s, err := goluaez.NewState("uu = function(a, b) return a == b end")
	c.Assert(err, IsNil)
	v, e := s.Run("uu(2, 9)")
	c.Assert(e, IsNil)
	c.Assert(v, Equals, false)

	v, e = s.Run("uu('aa', 'aa')")
	c.Assert(e, IsNil)
	c.Assert(v, Equals, true)
}

func (t *Tester) TestFuncNil(c *C) {
	s, err := goluaez.NewState("uu = function(a) end")
	c.Assert(err, IsNil)
	v, e := s.Run("uu(22)")
	c.Assert(e, IsNil)
	c.Assert(v, IsNil)

}

func (t *Tester) TestFuncTable(c *C) {
	s, err := goluaez.NewState(`uu = function() 
	    a = {}
		a["a"] = 1
		a["b"] = 2
		a["c"] = 3
		return a
	end`)
	c.Assert(err, IsNil)
	v, e := s.Run("uu()")
	c.Assert(e, IsNil)
	vv, ok := v.(map[string]interface{})
	c.Assert(ok, Equals, true)
	c.Assert(len(vv), Equals, 3)
	c.Assert(vv["b"], Equals, float64(2))
}

func (t *Tester) TestFuncSlice(c *C) {
	s, err := goluaez.NewState(`uu = function()
	    local a = {22,33,44,55}
		return a
	end`)
	c.Assert(err, IsNil)
	v, e := s.Run("uu()")
	c.Assert(e, IsNil)

	vv, ok := v.([]interface{})
	c.Assert(ok, Equals, true)
	c.Assert(len(vv), Equals, 4)
	c.Assert(vv[1], Equals, float64(33))

}

func (t *Tester) TestRun(c *C) {
	s, err := goluaez.NewState(`fn = function(a, b, c)
		return a + b + c
	end`)
	c.Assert(err, IsNil)

	v, err := s.Run("fn(a, b, c)", map[string]interface{}{"a": 22, "b": float64(22), "c": "22"})
	c.Assert(err, IsNil)

	c.Assert(v, Equals, float64(66))

}

func (t *Tester) TestSliceToTable(c *C) {
	s, _ := goluaez.NewState()

	v, err := s.Run("#tbl", map[string]interface{}{"tbl": []string{"aaa", "bbb", "cccc"}})
	c.Assert(err, IsNil)
	c.Assert(v, Equals, float64(3))

	v, err = s.Run("table.concat(tbl, ',')", map[string]interface{}{"tbl": []string{"aaa", "bbb", "cccc"}})
	c.Assert(err, IsNil)
	c.Assert(v, Equals, "aaa,bbb,cccc")

	/*
		v, err = s.Run("table.concat(tbl, ',')", map[string]interface{}{"tbl": [3]string{"aaa", "bbb", "cccc"}})
		c.Assert(err, IsNil)
		c.Assert(v, Equals, "aaa,bbb,cccc")
	*/

}

func (t *Tester) TestSplit(c *C) {
	s, _ := goluaez.NewState()

	v, err := s.Run("x:split('%s')", map[string]interface{}{"x": "a b c"})
	c.Assert(err, IsNil)
	c.Assert(len(v.([]interface{})), Equals, 3)
}

func (t *Tester) TestStrip(c *C) {
	s, _ := goluaez.NewState()

	v, err := s.Run("x:strip()", map[string]interface{}{"x": " a b c "})
	c.Assert(err, IsNil)
	c.Assert(v.(string), Equals, "a b c")
}

func (t *Tester) TestGo2LValue(c *C) {

	v, err := goluaez.Go2LValue(64)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, lua.LNumber(64))

	v, err = goluaez.Go2LValue("string")
	c.Assert(err, IsNil)
	c.Assert(v, Equals, lua.LString("string"))

	v, err = goluaez.Go2LValue([]string{"thing 1", "thing 2"})
	log.Println(v)
	c.Assert(err, IsNil)
	c.Assert(v.(*lua.LTable).Len(), Equals, 2)
	c.Assert(v.(*lua.LTable).RawGetInt(2).String(), Equals, "thing 2")

}
