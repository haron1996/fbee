package funcs

import (
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func LoginToFacebook() {
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	page := browser.MustPage("https://web.facebook.com/").MustWaitLoad()

	defer browser.MustClose()

	pageHasLoginForm, _, err := page.Has(`form[data-testid="royal_login_form"]`)
	if err != nil {
		log.Println("Error checking if page has login form:", err)
		return
	}
	if pageHasLoginForm {
		page.MustElement("#email").MustInput("hkibetr@gmail.com")
		page.MustElement("#pass").MustInput("33608080")
		page.MustElement("button").MustClick()

		time.Sleep(10 * time.Second)
	}

	log.Println("Login complete")
}
