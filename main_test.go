package main

import (
	"fmt"
	"testing"
)

func TestStuff(t *testing.T) {
	a := "10S"
	b := "0D"
	c := "1D"
	d := "2D"
	e := "JD"
	f := "QD"
	g := "KD"
	h := "AD"

	fmt.Println("10", a[0], 10)
	fmt.Println("0", b[0])
	fmt.Println("1", c[0])
	fmt.Println("2", d[0])
	fmt.Println("J", e[0])
	fmt.Println("Q", f[0])
	fmt.Println("K", g[0])
	fmt.Println("A", h[0])
}
