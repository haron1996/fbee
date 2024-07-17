// package main

// import (
// 	"time"

// 	"github.com/go-rod/rod"
// 	"github.com/go-rod/rod/lib/launcher"
// )

// func main() {
// 	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()

// 	browser := rod.New().ControlURL(u).MustConnect()

// 	page := browser.MustPage("https://web.facebook.com/").MustWaitLoad()

// 	defer browser.MustClose()

// 	page.MustElement("#email").MustInput("0718448461")
// 	page.MustElement("#pass").MustInput("33608080")
// 	page.MustElement("button").MustClick()

// 	time.Sleep(10 * time.Second)

// 	page = browser.MustPage("https://web.facebook.com/").MustWaitLoad()

// 	page.MustScreenshot("facebook.png")
// }

package main

import (
	"time"

	"github.com/haron1996/fb/funcs"
)

func main() {
	funcs.LoginToFacebook()
	time.Sleep(10 * time.Second)
	//funcs.PostToMarketplace()
	funcs.ListInMorePlaces()
	//funcs.PostToGroups()
	//funcs.Test()
}
