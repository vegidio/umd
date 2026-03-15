package main

import (
	"fmt"

	"github.com/vegidio/go-sak/fetch"
)

func main() {
	co := fetch.GetBrowserCookies("simpcity.cr")
	fmt.Printf("%s\n", co)
}
