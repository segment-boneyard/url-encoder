package encoder_test

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/segmentio/url-encoder"
)

func TestMarshalsString(t *testing.T) {
	d := struct {
		String string
	}{"foo"}

	v := encoder.Marshal(d)

	assert.Equal(t, "foo", v.Get("String"))
}

func TestMarshalsInt(t *testing.T) {
	d := struct {
		Int int
	}{9}

	v := encoder.Marshal(d)

	assert.Equal(t, "9", v.Get("Int"))
}

type A struct {
	Elem1 string
	Elem2 B
}

type B struct {
	Elem1 string
	Elem2 map[string]string
}

func TestMarshalsNestedStructs(t *testing.T) {
	d := A{
		Elem1: "foo",
		Elem2: B{
			Elem1: "bar",
			Elem2: map[string]string{"qaz": "qux"},
		},
	}

	v := encoder.Marshal(d)

	fmt.Println(v)

	assert.Equal(t, "foo", v.Get("Elem1"))
	assert.Equal(t, "bar", v.Get("Elem2.Elem1"))
	assert.Equal(t, "qux", v.Get("Elem2.Elem2.qaz"))
}

func TestMarshalsNestedMap(t *testing.T) {
	d := struct {
		Map map[string]int
	}{map[string]int{"foo": 4}}

	v := encoder.Marshal(d)

	assert.Equal(t, "4", v.Get("Map.foo"))
}

func TestMarshalsTopLevelMap(t *testing.T) {
	d := map[string]string{"foo": "4", "bar": "qaz"}

	v := encoder.Marshal(d)

	assert.Equal(t, "4", v.Get("foo"))
	assert.Equal(t, "qaz", v.Get("bar"))
}

func TestMarshalsSlices(t *testing.T) {
	d := struct {
		IntArray []int
	}{[]int{4, 5, 3}}

	v := encoder.Marshal(d)

	assert.Equal(t, []string{"4", "5", "3"}, v["IntArray"])
}

func intPtr(x int) *int {
	return &x
}

func stringPtr(x string) *string {
	return &x
}

func TestMarshalsPointers(t *testing.T) {
	d := struct {
		String *string
		Int    *int
	}{stringPtr("foo"), intPtr(3)}

	v := encoder.Marshal(d)

	assert.Equal(t, "foo", v.Get("String"))
	assert.Equal(t, "3", v.Get("Int"))
}

func TestMarshalsNilPointers(t *testing.T) {
	d := struct {
		String *string `url:"string_ptr"`
		Int    *int
	}{nil, nil}

	v := encoder.Marshal(d)

	if _, ok := v["string_ptr"]; ok {
		t.Errorf("expected nothing for String pointer but got %q", v["string_ptr"])
	}
	if _, ok := v["Int"]; ok {
		t.Errorf("expected nothing for IntPtr but got %q", v["Int"])
	}
}

func TestRespectsCustomKeyValues(t *testing.T) {
	d := struct {
		String string `url:"bar"`
	}{"foo"}

	v := encoder.Marshal(d)

	assert.Equal(t, "foo", v.Get("bar"))
	assert.Equal(t, "", v.Get("String"))
}
