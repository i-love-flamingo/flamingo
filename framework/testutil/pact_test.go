package testutil

import (
	"encoding/json"
	"testing"
)

type teststruct struct {
	Foo  string      `json:"foo"`
	Bar  int         `json:"bar"`
	Blub *teststruct `json:"blub"`
}

func TestPactEncodeLike(t *testing.T) {
	t.Skip("not implemented yet")

	var test = teststruct{
		Foo: "Foo String",
		Bar: 100,
		Blub: &teststruct{
			Foo: "String2",
		},
	}

	b, _ := json.Marshal(test)
	var testencoded = PactEncodeLike(b)

	if testencoded != `` {
		//t.Fatal("wrong encoding", testencoded)
	}
}
