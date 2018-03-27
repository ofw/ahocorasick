## ⚡ Aho-Corasick Pattern Matching Algorithm ⚡

Aho-Corasick string matching algorithm for golang.

This if a fork of original https://github.com/gansidui/ahocorasick library which is not updated since 2014.

Key improvements:
* Thread safety for multiple calls to `Match` method 🌪️
* Perfomance optimizations (about 5x speed and reduced allocations) 🏎:
```
    BenchmarkOriginal      1000000          1424 ns/op          64 B/op          4 allocs/op
    BenchmarkOptimized     5000000           237 ns/op          35 B/op          2 allocs/op
```
* Fixed incorrect results with some test cases

Now this package is even faster than https://github.com/cloudflare/ahocorasick:

```
BenchmarkOfw-4          	 5000000	       318 ns/op	      16 B/op	       2 allocs/op
BenchmarkCloudflare-4   	 3000000	       455 ns/op	     104 B/op	       4 allocs/op
```

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