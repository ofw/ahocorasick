## ⚡ Aho-Corasick Pattern Matching Algorithm ⚡

Aho-Corasick string matching algorithm for golang.

This if a fork of original https://github.com/gansidui/ahocorasick library which is not updated since 2014.

Key improvements:
* Thread safety for multiple calls to `Match` method 🌪️
* Perfomance optimizations (about 5x speed and reduced allocations) 🏎
* Fixed incorrect results with some test cases


~~~ go
package main

import (
	"fmt"
	"github.com/ofw/ahocorasick"
)

func main() {
	ac := ahocorasick.NewMatcher()

	dictionary := []string{"hello", "world", "世界", "google", "golang", "c++", "love"}

	ac.Build(dictionary)

	ret := ac.Match("hello世界, hello google, i love golang!!!")

	for _, i := range ret {
		fmt.Println(dictionary[i])
	}
}


~~~

## License

MIT