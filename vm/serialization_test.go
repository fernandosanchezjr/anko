package vm

import (
	"bytes"
	"encoding/gob"
	"github.com/fernandosanchezjr/anko/ast"
	"github.com/fernandosanchezjr/anko/parser"
	"github.com/kylelemons/godebug/pretty"
	"testing"
)

var source []string = []string{
	`a = nil`,
	`a = 1`,
	`a = 1.2`,
	`a = "foo"`,
	`a = true`,
	`a = false`,
	`a = [1,2,3]`,
	`a = {"foo": "bar", "bar": "baz"}`,
	`a = {"foo": "bar", "bar": {"blah": true, "blah!": [1.3e3, true]}}`,
	`toString(1)`,
	`toString(1.2)`,
	`toString(true)`,
	`toString(false)`,
	`toString("foo")`,
	`1 > 0`,
	`1 == 1.0`,
	`1 != "1"`,
	`1 != "1"`,
	`1.1 == 1.1`,
	`"1" == "1"`,
	`false != "1"`,
	`false != true`,
	`false == false`,
	`true == true`,
	`nil == nil`,
	`1 <= 1`,
	`1.0 <= 1.0`,
	`1 <= 2 ? true : false`,
	`a = 1; a += 1`,
	`a = 2; a -= 1`,
	`a = 2; a *= 2`,
	`a = 3; a /= 2`,
	`a = 2; a++`,
	`a = 2; a--`,
	`a = 2**3`,
	`a = 1; a &= 2`,
	`a = 1; a |= 2`,
	`a = !3`,
	`a = !true`,
	`a = !false`,
	`a = ^3`,
	`a = 3 << 2`,
	`a = 11 >> 2`,
	`func a() { return 2 }`,
	`func b(x) { return x + 1 }`,
	`func c(x) { return x, x + 1 }`,
	`func d(x) { return func() { return x + 1 } }`,
	`var x = func(x) {
	  return func(y) {
	    x(y)
	  }
	 }(func(z) {
	  return "Yay! " + z
	 })("hello world")`,
	`len("foo")`,
	`len("")`,
	`len([1,2,true,["foo"]])`,
	`x = 0
     for a in [1,2,3] {
      x += 1
     }`,
	`x = 0
	 for {
	  x += 1
	  if (x > 3) {
	    break
	  }
	 }`,
	`func loop_with_return_stmt() {
	  y = 0
	  for {
	    if y == 5 {
	      return y
	    }
	    y++
	  }
	  return 1
	 }`,
	`func for_with_return_stmt() {
	  y = 0
	  for k in range(0, 10) {
	    if k == 5 {
	      return y
	    }
	    y++
	  }
	  return 1
	 }`,
	`x = 0
	 for a = 0; a < 10; a++ {
	  x++
	 }`,
	`func cstylefor_with_return_stmt() {
	  y = 0
	  for i = 0; i < 10; i++ {
	    if i == 5 {
	      return y
	    }
	    y++
	  }
	  return 1
	 }`,
	`resp = {
	    "items": [{
	        "someData": 2,
	    }]
	 }
	 x = 0
	 for item in resp.items {
	    x += item.someData
	 }`,
	`x = 0
	 r = -1
	 switch x {
	 case 0:
	  r = 0
	 case 1:
	  r = 1
	 case 2:
	  r = 2
	 }`,
	`x = 3
	 r = -1
	 switch x {
	 case 0:
	  r = 0
	 case 1:
	  r = 1
	 case 2:
	  r = 2
	 }`,
	`x = 3
	 r = -1
	 switch x {
	 case 0:
	  r = 0
	 case 1:
	  r = 1
	 case 2:
	  r = 2
	 default:
	  r = 3
	 }`,
	`r = -1
	 if (false) {
	  r = 1
	 } else if (false) {
	  r = 2
	 } else if (false) {
	  r = 3
	 } else {
	  r = 4
	 }`,
	`a = toByteSlice("あいうえお")
	 b = [227, 129, 130, 227, 129, 132, 227, 129, 134, 227, 129, 136, 227, 129, 138]
	 x = 0
	 for i = 0; i < len(a); i++ {
	  if (a[i] == b[i]) {
	    x++
	  }
	 }`,
	`a = toRuneSlice("あいうえお")
	 b = [12354, 12356, 12358, 12360, 12362]
	 x = 0
	 for i = 0; i < len(a); i++ {
	  if (a[i] == b[i]) {
	    x++
	  }
	 }`}

func TestParsedSourceSerialization(t *testing.T) {
	for _, value := range source {
		st, err := parser.ParseSrc(value)
		if err != nil {
			t.Fatal(err)
		}
		dst := make([]ast.Stmt, 0)
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(st); err != nil {
			t.Fatal("Error encoding", value, " - ", err)
		}
		decoder := gob.NewDecoder(&buf)
		if err := decoder.Decode(&dst); err != nil {
			t.Fatal("Error decoding", value, " - ", err)
		}
		if diff := pretty.Compare(dst, st); len(diff) > 0 {
			t.Fatal("Encode/decode mismatch with", value, "\n", diff)
		}
	}
}
