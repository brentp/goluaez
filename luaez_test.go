package goluaez_test

import (
	"testing"

	"github.com/brentp/goluaez"

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

func (t *Tester) TestSplit(c *C) {
	s, err := goluaez.NewState()
	c.Assert(err, IsNil)
	res, err := s.Run("split('xexexex', 'e')")
	c.Assert(err, IsNil)
	c.Assert(len(res.([]string)), Equals, 4)

	res, err = s.Run("split('xe22efg22ex', '\\\\d+')")
	c.Assert(err, IsNil)
	c.Assert(len(res.([]string)), Equals, 3)
	c.Assert(res.([]string)[1], Equals, "efg")

}

func (t *Tester) TestIndex(c *C) {
	s, err := goluaez.NewState()
	c.Assert(err, IsNil)
	res, err := s.Run("index('abcdefg', 'c')")
	c.Assert(err, IsNil)
	c.Assert(int(res.(float64)), Equals, 2)

	res, err = s.Run("index('abcdefg', 'h')")
	c.Assert(err, IsNil)
	c.Assert(int(res.(float64)), Equals, -1)
}
