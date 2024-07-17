package funcs

import (
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/haron1996/fb/config"
)

func LoginToFacebook() {
	// load config files
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}
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
		page.MustElement("#email").MustInput(config.Username)
		page.MustElement("#pass").MustInput(config.Password)
		page.MustElement("button").MustClick()
	}
	time.Sleep(10 * time.Second)
	page.MustScreenshot("home.png")
	log.Println("Login complete")
}
