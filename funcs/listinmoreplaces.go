package funcs

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func ListInMorePlaces() {

	// Launch the browser with specific configurations
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // Ensure the browser closes when main function exits

	page := browser.MustPage("https://web.facebook.com/marketplace/you/selling").MustWaitLoad()

	time.Sleep(10 * time.Second)

	menus := page.MustElements(`div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.x1k70j0n.xzueoph.xzboxd6.x14l7nz5 > div.x78zum5.x1n2onr6.xh8yej3 > div.x9f619.x1n2onr6.x1ja2u2z.x1jx94hy.x1qpq9i9.xdney7k.xu5ydu1.xt3gfkd.xh8yej3.x6ikm8r.x10wlt62.xquyuld`)

	for _, menu := range menus {
		menu.MustClick()
		time.Sleep(10 * time.Second)
		pageHasInfo, info, err := page.Has(`div[aria-label="Your listing"]`)
		if err != nil {
			log.Println("Error checking if page has info:", err)
			return
		}
		if !pageHasInfo {
			log.Println("Page has no info")
			return
		}
		list := info.MustElement(`div[aria-label="List to more places"]`)
		list.MustClick()
		time.Sleep(10 * time.Second)
		pageHasCard, card, err := page.Has(`div[aria-label="List in more places"]`)
		if err != nil {
			log.Println("Error checking if page has list in more places card:", err)
			return
		}
		if !pageHasCard {
			log.Println("Page has no list in more places card")
			return
		}
		containers := card.MustElements("div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.x1e558r4.x150jy0e")
		container := containers[1]
		groups := container.MustElements(`div[data-visualcompletion="ignore-dynamic"][style="padding-left: 8px; padding-right: 8px;"]`)
		for _, group := range groups {
			group.MustClick()
		}
		fmt.Printf("Total Groups:%d\n", len(groups))
		time.Sleep(5 * time.Second)
		card.MustElements("div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x193iq5w.xeuugli.x1iyjqo2.xs83m0k.x150jy0e.x1e558r4.xjkvuk6.x1iorvi4.xdl72j9")[1].MustElement(`div[aria-label="Post"]`).MustClick()
		time.Sleep(10 * time.Second)
		info.MustElement(`div[aria-label="Close"]`).MustClick()
		time.Sleep(10 * time.Second)
		log.Println("Listing shared to groups")
	}

	page.MustScreenshot("facebook.png")
}
