package dao

import (
	"fmt"
	"testing"
)

type T struct {
}

func Test(t *testing.T) {
	ts := make([]*T, 2)
	ts = append(ts, &T{})
	ts = append(ts, nil)
	fmt.Println(ts)
}
