package main

import (
	"time"

	"github.com/haron1996/fb/funcs"
)

func main() {
	funcs.LoginToFacebook()
	time.Sleep(5 * time.Second)
	funcs.PostToMarketplace()
	funcs.ListInMorePlaces()
}
