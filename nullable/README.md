# nullable

`nullable` makes it easy to handle JSON Strings, Ints, and Bools with values of `null`.

## Description

`nullable` is a library that wraps `string`, `bool`, `int`, and `int64` and provides custom `UnmarhsalJSON` functions that handle a null value by defaulting to the type's default value.

## Installation

    go get github.com/samuelkadolph/go/nullable

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	n "github.com/samuelkadolph/go/nullable"
)

type Widget struct {
	Name  n.String `json:"name"`
	Count n.Int    `json:"enabled"`
}

func main() {
	b := []byte(`{"name":null,"count":null}`)
	w := Widget{}

	json.Unmarshal(b, &w)

	fmt.Printf("'%s' %d\n", w.Name, w.Count) // '' 0

	if int(w.Count) == 0 {
		fmt.Printf("count is zero\n")
	}
}
```
