package funcs

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func ListInMorePlaces() {
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	defer browser.MustClose()

	page := browser.MustPage("https://web.facebook.com/marketplace/you/selling").MustWaitLoad()

	for i := 0; i < 30; i++ {
		err := page.Mouse.Scroll(0, 1000, 0)
		if err != nil {
			log.Println("Error scrolling page:", err)
			return
		}
		time.Sleep(15 * time.Second)
	}

	menus := page.MustElements(`div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.x1k70j0n.xzueoph.xzboxd6.x14l7nz5 > div.x78zum5.x1n2onr6.xh8yej3 > div.x9f619.x1n2onr6.x1ja2u2z.x1jx94hy.x1qpq9i9.xdney7k.xu5ydu1.xt3gfkd.xh8yej3.x6ikm8r.x10wlt62.xquyuld`)

	totalMenus := len(menus)

	fmt.Println("Total returned menus:", totalMenus)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(totalMenus, func(i, j int) {
		menus[i], menus[j] = menus[j], menus[i]
	})

	// Select the first 20 elements after shuffling
	randomMenus := menus[:20]

	for _, menu := range randomMenus {
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

		infoCardHasListBtn, listBtn, err := info.Has(`div[aria-label="List to more places"]`)
		if err != nil {
			log.Println("Error checking if info card has list in more places button:", err)
			return
		}
		if !infoCardHasListBtn {
			fmt.Println("Listing has no list in more places button")
			info.MustElement(`div[aria-label="Close"]`).MustClick()
			continue
		}
		listBtn.MustClick()
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
		time.Sleep(10 * time.Second)
		card.MustElements("div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x193iq5w.xeuugli.x1iyjqo2.xs83m0k.x150jy0e.x1e558r4.xjkvuk6.x1iorvi4.xdl72j9")[1].MustElement(`div[aria-label="Post"]`).MustClick()
		time.Sleep(10 * time.Second)
		info.MustElement(`div[aria-label="Close"]`).MustClick()
		time.Sleep(10 * time.Second)
		fmt.Printf("Listing shared to %d groups\n", len(groups))
	}

	page.MustScreenshot("facebook.png")
}
