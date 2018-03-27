package ahocorasick

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var assert = require.True

func Test1(t *testing.T) {
	ac := NewMatcher()

	dictionary := []string{"she", "he", "say", "shr", "her"}
	ac.Build(dictionary)

	ret := ac.Match("yasherhs")
	if len(ret) != 3 || ret[0] != 0 || ret[1] != 1 || ret[2] != 4 {
		t.Fatal()
	}

	ret = ac.Match("yasherhs")
	if len(ret) != 3 || ret[0] != 0 || ret[1] != 1 || ret[2] != 4 {
		t.Fatal()
	}

	if ac.GetMatchResultSize("yasherhs") != 3 {
		t.Fatal()
	}
}

func Test2(t *testing.T) {
	ac := NewMatcher()

	dictionary := []string{"hello", "世界", "hello世界", "hello"}
	ac.Build(dictionary)

	ret := ac.Match("hello世界")
	if len(ret) != 4 {
		t.Fatal()
	}

	ret = ac.Match("世界")
	if len(ret) != 1 {
		t.Fatal()
	}

	ret = ac.Match("hello")
	if len(ret) != 2 {
		t.Fatal()
	}
}

func Test3(t *testing.T) {
	ac := NewMatcher()

	dictionary := []string{"abc", "bc", "ac", "bc", "de", "efg", "fgh", "hi", "abcd", "ac"}
	ac.Build(dictionary)

	ret := ac.Match("abcdefghij")
	if len(ret) != ac.GetMatchResultSize("abcdefghij") || len(ret) != 8 {
		t.Fatal()
	}

	ret = ac.Match("abcdef")
	if len(ret) != 5 {
		t.Fatal()
	}

	ret = ac.Match("acdejefg")
	if len(ret) != 4 {
		t.Fatal()
	}

	if len(ac.Match("abcd")) != 4 {
		t.Fatal()
	}

	if len(ac.Match("adefacde")) != 3 {
		t.Fatal()
	}

	ret = ac.Match("agbdfgiadafgha")
	if len(ret) != 1 || dictionary[ret[0]] != "fgh" {
		t.Fatal()
	}
}

func TestConcurrent(t *testing.T) {
	ac := NewMatcher()
	ac.Build([]string{"foo", "bфr", "baz"})
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			require.Equal(t, []uint32{0, 1, 2}, ac.Match("foobфrbaz"))
		}()
	}
	wg.Wait()
}

func BenchmarkOfw(b *testing.B) {

	ac := NewMatcher()
	ac.Build([]string{"foo", "bar", "baz"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ac.Match("fooasldkjflaksjbarsdfasdfbazasdfdf")
	}
	b.ReportAllocs()
}

//func BenchmarkCloudflare(b *testing.B) {
//	ac := ahocorasick.NewStringMatcher([]string{"foo", "bar", "baz"})
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		ac.Match([]byte("fooasldkjflaksjbarsdfasdfbazasdfdf"))
//	}
//}

func TestNoPatterns(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{})
	hits := m.Match("foo bar baz")
	assert(t, len(hits) == 0)
}

func TestNoData(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"foo", "baz", "bar"})
	hits := m.Match("")
	assert(t, len(hits) == 0)
}

func TestSuffixes(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"Superman", "uperman", "perman", "erman"})
	hits := m.Match(string("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 3)
}

func TestPrefixes(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"Superman", "Superma", "Superm", "Super"})
	hits := m.Match(string("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 3)
	assert(t, hits[1] == 2)
	assert(t, hits[2] == 1)
	assert(t, hits[3] == 0)
}

func TestInterior(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"Steel", "tee", "e"})
	hits := m.Match(string("The Man Of Steel: Superman"))
	fmt.Println(hits)
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 2)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 0)
}

func TestMatchAtStart(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"The", "Th", "he"})
	hits := m.Match(string("The Man Of Steel: Superman"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 0)
	assert(t, hits[2] == 2)
}

func TestMatchAtEnd(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"teel", "eel", "el"})
	hits := m.Match(string("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
}

func TestOverlappingPatterns(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"Man ", "n Of", "Of S"})
	hits := m.Match(string("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
}

func TestMultipleMatches(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"The", "Man", "an"})
	hits := m.Match(string("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 2)
	assert(t, hits[2] == 0)
}

func TestSingleCharacterMatches(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"a", "M", "z"})
	hits := m.Match(string("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 1)
	assert(t, hits[1] == 0)
}

func TestNothingMatches(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"baz", "bar", "foo"})
	hits := m.Match(string("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 0)
}

func TestWikipedia(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"a", "ab", "bc", "bca", "c", "caa"})
	hits := m.Match(string("abccab"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 4)

	hits = m.Match(string("bccab"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 2)
	assert(t, hits[1] == 4)
	assert(t, hits[2] == 0)
	assert(t, hits[3] == 1)

	hits = m.Match(string("bccb"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 2)
	assert(t, hits[1] == 4)
}

func TestMatch(t *testing.T) {
	m := NewMatcher()
	m.Build([]string{"Mozilla", "Mac", "Macintosh", "Safari", "Sausage"})
	hits := m.Match(string("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 2)
	assert(t, hits[3] == 3)

	hits = m.Match(string("Mozilla/5.0 (Mac; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 1)
	assert(t, hits[2] == 3)

	hits = m.Match(string("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == 0)
	assert(t, hits[1] == 3)

	hits = m.Match(string("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 1)
	assert(t, hits[0] == 0)

	hits = m.Match(string("Mazilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 0)
}
