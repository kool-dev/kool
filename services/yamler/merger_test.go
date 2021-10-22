package yamler

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func testMerger(t *testing.T, foo, bar, expected string) {
	var y1, y2 = new(yaml.Node), new(yaml.Node)

	_ = yaml.Unmarshal([]byte(foo), y1)
	_ = yaml.Unmarshal([]byte(bar), y2)

	m := &DefaultMerger{}

	if err := m.Merge(y2, y1); err != nil {
		t.Errorf("unexpected error merging yaml: %v", err)
	}

	bs, _ := yaml.Marshal(y1)
	merged := strings.Trim(string(bs), " \t\n")
	if merged != expected {
		t.Errorf("unexpected merge result: expected '%s' but got '%s'", expected, merged)
	}
}

func TestMergerMap(t *testing.T) {
	foo := "foo: yyy"
	bar := "bar: xxx"
	expected := foo + "\n" + bar

	testMerger(t, foo, bar, expected)
}

func TestMergerMapOverride(t *testing.T) {
	foo := "foo: yyy"
	bar := "bar: xxx\nfoo: zzz"
	expected := "foo: zzz\nbar: xxx"

	testMerger(t, foo, bar, expected)
}

func TestMergerComments(t *testing.T) {
	comment := "# comment\nfoo: test"
	bar := "bar: xxx"
	expected := comment + "\n" + bar

	testMerger(t, comment, bar, expected)
}

func TestMergerListAppend(t *testing.T) {
	foo := "foo: [yyy, zzz]"
	bar := "foo: [xxx]"
	expected := "foo: [yyy, zzz, xxx]"

	testMerger(t, foo, bar, expected)
}

func TestMergerMapNested(t *testing.T) {
	foo := "foo:"
	bar := "foo:\n  bar: xxx"
	expected := "foo:\n    bar: xxx"

	testMerger(t, foo, bar, expected)
}

func TestMergerMapNestedOverride(t *testing.T) {
	foo := "foo:\n  bar: yyy\n  zuk: ppp"
	bar := "foo:\n  bar: xxx"
	expected := "foo:\n    bar: xxx\n    zuk: ppp"

	testMerger(t, foo, bar, expected)
}
