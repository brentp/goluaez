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
