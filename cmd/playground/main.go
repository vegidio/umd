package main

import (
	"github.com/vegidio/go-sak/fetch"
)

func main() {
	co := fetch.GetBrowserCookies("simpcity.cr")
	println("%s", co)
}
