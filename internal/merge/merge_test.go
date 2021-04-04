package merge

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sno6/gosane/internal/types"
)

type InternalSimpleTest struct {
	Key   string
	Value string
}

type InternalEmbededTest struct {
	Key                string
	Value              string
	InternalSimpleTest InternalSimpleTest
}

func TestMergerMergeSimple(t *testing.T) {
	src := InternalSimpleTest{
		Key:   "Key-src",
		Value: "Value-src",
	}

	// Should get the Key and Value written to it
	dst := InternalSimpleTest{
		Key: "Key-dst",
	}

	err := Merge(&dst, src)
	if err != nil {
		t.Fail()
	}

	if !strings.Contains(dst.Key, "src") && !strings.Contains(dst.Value, "src") {
		t.Fail()
	}
}

func TestMergerMergeEmbeded(t *testing.T) {
	src := InternalEmbededTest{
		Key:   "Key-outer-src",
		Value: "Value-outer-src",
		InternalSimpleTest: InternalSimpleTest{
			Key: "Key-inner-src",
		},
	}

	// Should get the Key and Value written to it
	dst := InternalEmbededTest{
		Key: "Key-dst",
		InternalSimpleTest: InternalSimpleTest{
			Value: "Key-inner-dst",
		},
	}

	err := Merge(&dst, src)
	if err != nil {
		t.Fail()
	}

	if !strings.Contains(dst.InternalSimpleTest.Key, "src") &&
		!strings.Contains(dst.InternalSimpleTest.Value, "src") {
		t.Fail()
	}
}

func TestMergeMapInit(t *testing.T) {
	b := `{"key":"value"}`

	var m Map
	if err := m.Init([]byte(b)); err != nil {
		t.Error(err)
		return
	}

	v, ok := m["key"]
	if !ok || v != "value" {
		t.Fail()
		return
	}
}

func TestMergeMapMerge(t *testing.T) {
	type CustomType string

	type S struct {
		Str    string     `json:"str,omitempty"`
		StrPtr *string    `json:"str_ptr,omitempty"`
		Custom CustomType `json:"custom,omitempty"`
	}

	type S2 struct {
		Str    string     `json:"str,omitempty"`
		StrPtr *string    `json:"strPtr,omitempty"`
		Custom CustomType `json:"custom,omitempty"`
	}

	existing := S{Str: "Test"}
	updated := S2{Str: "Updated-String", StrPtr: types.String("Updated-Pointer"), Custom: CustomType("Custom-Type")}

	b, err := json.Marshal(updated)
	if err != nil {
		t.Error(err)
		return
	}

	var m Map
	if err := m.Init(b); err != nil {
		t.Error(err)
		return
	}

	if err := MergeMap(&existing, updated, m, true); err != nil {
		t.Error(err)
		return
	}

	if existing.Str != updated.Str || *existing.StrPtr != *updated.StrPtr || existing.Custom != updated.Custom {
		t.Log("Values don't match")
		t.Fail()
		return
	}
}
