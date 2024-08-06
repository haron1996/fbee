package funcs

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/haron1996/fb/config"
)

func CreatePost() {
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	defer browser.MustClose()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	root := config.Root
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading root directory:", err)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	for _, entry := range entries {
		page := browser.MustPage("https://web.facebook.com/profile.php?id=61563372774438").MustWaitLoad()

		time.Sleep(10 * time.Second)

		page.MustScreenshot("fb.png")

		btns, err := page.Elements(`div.x1i10hfl.xjbqb8w.xjqpnuy.xa49m3k.xqeqjp1.x2hbi6w.x13fuv20.xu3j5b3.x1q0q8m5.x26u7qi.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x2lwn1j.x1n2onr6.x16tdsg8.x1hl2dhg.xggy1nq.x1ja2u2z.x1t137rt.x1q0g3np.x87ps6o.x1lku1pv.x1a2a7pz.x6s0dn4.x1lq5wgf.xgqcy7u.x30kzoy.x9jhf4c.x78zum5.x1r8uery.x1iyjqo2.xs83m0k.xl56j7k.x1pshirs.x1y1aw1k.x1sxyh0.xwib8y2.xurb0ha`)
		if err != nil {
			log.Println("Error getting elements:", err)
			return
		}

		photoVideoBtn := btns[1]

		err = photoVideoBtn.Click("left", 1)
		if err != nil {
			log.Println("Error clicking photo/video button:", err)
			return
		}

		time.Sleep(10 * time.Second)

		pageHasDialog, dialog, err := page.Has(`div[role="dialog"]`)
		if err != nil {
			log.Println("Error checkin if pages has dialog:", err)
			return
		}

		if !pageHasDialog {
			log.Println("Page ha no dialog:", err)
			return
		}

		dialogHasFileInput, fileInput, err := dialog.Has(`input[type="file"]`)
		if err != nil {
			log.Println("Error checking if dialog has file input:", err)
			return
		}

		if !dialogHasFileInput {
			log.Println("Dialog has no file input")
			return
		}

		subDir := filepath.Join(root, entry.Name())
		fmt.Println("DIRECTORY:", subDir)

		subEntries, err := os.ReadDir(subDir)
		if err != nil {
			fmt.Println("Error reading subdirectory:", err)
			continue
		}

		var imageFiles []string

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				filePath := filepath.Join(subDir, subEntry.Name())
				if filepath.Ext(filePath) != ".txt" {
					fmt.Println("IMAGE:", filePath)
					imageFiles = append(imageFiles, filePath)
				}
			}
		}

		if len(imageFiles) > 0 {
			fileInput.MustSetFiles(imageFiles...)
		}

		detailsFile := filepath.Join(subDir, "details.txt")
		file, err := os.Open(detailsFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer file.Close()

		var title, description string

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "title:"):
				title = strings.TrimSpace(line[len("title:"):])
			case strings.HasPrefix(line, "description:"):
				description = strings.TrimSpace(line[len("description:"):])
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		parts := strings.Split(description, "...")

		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		formattedDescription := strings.Join(parts, "\n\n")

		fmt.Println("Description: " + formattedDescription)

		contentEditableDiv, err := dialog.Element(`div[aria-label="What's on your mind?"]`)
		if err != nil {
			log.Println("Error getting content editable div:", err)
			return
		}

		err = contentEditableDiv.Input(formattedDescription)
		if err != nil {
			log.Println("Error inputing phone description:", err)
			return
		}

		postBtn, err := dialog.Element(`div[aria-label="Post"]`)
		if err != nil {
			log.Println("Error getting post button:", err)
			return
		}

		time.Sleep(30 * time.Second)
		postBtn.MustScreenshot("btn.png")

		err = postBtn.Click("left", 1)
		if err != nil {
			log.Println("Error clicking post button:", err)
			return
		}

		time.Sleep(30 * time.Second)
		page.MustScreenshot("fb.png")

		fmt.Printf("%s posted successfully\n", title)
	}
}
