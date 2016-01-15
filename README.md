luaez
=====

[![GoDoc](https://godoc.org/github.com/brentp/goluaez?status.svg)](https://godoc.org/github.com/brentp/goluaez)
[![Build Status](https://travis-ci.org/brentp/goluaez.svg)](https://travis-ci.org/brentp/goluaez)
[![Coverage Status](https://coveralls.io/repos/brentp/goluaez/badge.svg?branch=master&service=github)](https://coveralls.io/github/brentp/goluaez?branch=master)


Easy embedding of lua in go.

`goluaez` wraps [gopher-lua](https://github.com/yuin/gopher-lua) and [gopher-luar](https://github.com/layeh/gopher-luar). `gopher-luar` does a nice job of converting go values to lua. This package uses that and also converts lua values to go and converts go slices to lua tables.

This makes it easy to have small 1 or 2 line user-defined functions in a go application that embeds lua.

Example usage:

```Go
package main

import (
	"fmt"
	"log"

	"github.com/brentp/goluaez"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	L, err := goluaez.NewState(`
adder = function(a, b)
    return a + b
end`)
	check(err)
	defer L.Close()

	// Run uses a mutex so can be run in a goroutine.
	result, err := L.Run("adder(x, y)", map[string]interface{}{"x": 12, "y": "23"})
	check(err)
	fmt.Println(result.(float64))
	// 35

	result, err = L.Run("'hello' .. ' world'")
	check(err)
	fmt.Println(result.(string))
	// hello world

	L.DoString("a = {}; a['a'] = 22; a['b'] = 33;")
	result, err = L.Run("a")
	check(err)
	fmt.Println(result.(map[string]interface{}))
	// map[b:33 a:22]

}
```

Prelude
-------

The embedded lua engine is populated with the lua functions in data/prelude.lua
So far, this consists of:

+ string:split(sep) (split a string by a delim)
+ string:strip() (remove whitespace from a string)
