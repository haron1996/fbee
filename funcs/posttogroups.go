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
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/haron1996/fb/config"
)

func PostToGroups() {
	// Launch the browser with specific configurations
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // Ensure the browser closes when main function exits

	page := browser.MustPage("https://web.facebook.com/groups/joins").MustWaitLoad()

	// count := 1000

	// for i := 0; i < count; i++ {
	// 	err := page.Mouse.Scroll(0, 1000000, 0)
	// 	if err != nil {
	// 		log.Println("Error scrolling page:", err)
	// 		return
	// 	}
	// }

	// get groups joined container
	classSelector := ".x9f619.x1gryazu.xkrivgy.x1ikqzku.x1h0ha7o.xg83lxy.xh8yej3"

	parents := page.MustElements(classSelector)

	var parent *rod.Element

	parentsLength := len(parents)

	log.Println(parentsLength)

	if parentsLength == 1 {
		parent = parents[0]
	} else if parentsLength == 2 {
		parent = parents[1]
	} else {
		log.Println("no parents found")
		return
	}

	selector := ".x1i10hfl.x1qjc9v5.xjbqb8w.xjqpnuy.xa49m3k.xqeqjp1.x2hbi6w.x13fuv20.xu3j5b3.x1q0q8m5.x26u7qi.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xdl72j9.x2lah0s.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x2lwn1j.xeuugli.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1n2onr6.x16tdsg8.x1hl2dhg.xggy1nq.x1ja2u2z.x1t137rt.x1o1ewxj.x3x9cwd.x1e5q0jg.x13rtm0m.x1q0g3np.x87ps6o.x1lku1pv.x1rg5ohu.x1a2a7pz.xh8yej3"

	anchors := parent.MustElements(selector)

	for _, a := range anchors {
		href := *a.MustAttribute("href")
		fmt.Println(href)
		page = browser.MustPage(href).MustWaitLoad()
		// time.Sleep(10 * time.Second)
		// page.MustElement(`div[aria-label="Sell Something"]`).MustClick()
		// time.Sleep(10 * time.Second)
		// page.MustElements("span.x1lliihq.x1iyjqo2")[0].MustClick()
		// load config files
		config, err := config.LoadConfig(".")
		if err != nil {
			log.Println("Error loading config:", err)
			return
		}
		// Root directory containing subdirectories with images
		root := config.Root
		entries, err := os.ReadDir(root)
		if err != nil {
			fmt.Println("Error reading root directory:", err)
			return
		}
		// Shuffle the entries slice
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(entries), func(i, j int) {
			entries[i], entries[j] = entries[j], entries[i]
		})

		// get images selector
		// imagesSelector := page.MustElement("div.x1gslohp.x1swvt13.x1pi30zi")

		// card := page.MustElement("div.x9f619.x1ja2u2z.x1k90msu.x6o7n8i.x1qfuztq.x10l6tqk.x17qophe.x13vifvy.x1hc1fzr.x71s49j.xh8yej3")

		fmt.Println("reached here...")

		for _, entry := range entries {
			log.Println("started looping entries...")
			time.Sleep(10 * time.Second)
			pageHasSellBtn, sellBtn, err := page.Has(`div[aria-label="Sell Something"]`)
			if err != nil {
				log.Println("Error checking if page/group has sell button:", err)
				return
			}

			if !pageHasSellBtn {
				continue
			}
			// sellBtn, err := page.Element(`div[aria-label="Sell Something"]`)
			// if err != nil {
			// 	log.Println("Error getting sell button element:", err)
			// 	return
			// }
			err = sellBtn.Click("left", 1)
			if err != nil {
				log.Println("Error clicking sale button element:", err)
				return
			}
			fmt.Println("reached here...")
			time.Sleep(10 * time.Second)
			page.MustElements("span.x1lliihq.x1iyjqo2")[0].MustWaitLoad().MustClick()
			fmt.Println("reached here...")
			time.Sleep(10 * time.Second)
			page.MustScreenshot("start.png")
			imagesSelector := page.MustElement("div.x1gslohp.x1swvt13.x1pi30zi")
			fmt.Println("reached here...")
			card := page.MustElement("div.x9f619.x1ja2u2z.x1k90msu.x6o7n8i.x1qfuztq.x10l6tqk.x17qophe.x13vifvy.x1hc1fzr.x71s49j.xh8yej3")
			fmt.Println("reached here...?")
			// Locate the file input element on the page
			fileInput := imagesSelector.MustElement(`input[type="file"]`)
			fmt.Println("file input located....")

			//log.Println(fileInput.HTML())

			// Path to the current subdirectory
			subDir := filepath.Join(root, entry.Name())
			fmt.Println("DIRECTORY:", subDir)

			// Read files from the subdirectory
			subEntries, err := os.ReadDir(subDir)
			if err != nil {
				fmt.Println("Error reading subdirectory:", err)
				continue
			}

			var imageFiles []string
			// Loop through files in the subdirectory
			for _, subEntry := range subEntries {
				if !subEntry.IsDir() {
					filePath := filepath.Join(subDir, subEntry.Name())
					if filepath.Ext(filePath) != ".txt" { // Check if the file is not a .txt file
						fmt.Println("IMAGE:", filePath)

						// Collect file paths to attach
						imageFiles = append(imageFiles, filePath)
					}
				}
			}

			// Attach collected files to the file input element
			if len(imageFiles) > 0 {
				fileInput.MustSetFiles(imageFiles...)
			}

			// Open the details.txt file within the subdirectory
			detailsFile := filepath.Join(subDir, "details.txt")
			file, err := os.Open(detailsFile)
			if err != nil {
				fmt.Println("Error opening file:", err)
				continue
			}
			defer file.Close()

			// Initialize variables to hold the extracted fields
			var title, price, description string

			// Create a new scanner to read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				switch {
				case strings.HasPrefix(line, "title:"):
					title = strings.TrimSpace(line[len("title:"):])
				case strings.HasPrefix(line, "price:"):
					price = strings.TrimSpace(line[len("price:"):])
				case strings.HasPrefix(line, "description:"):
					description = strings.TrimSpace(line[len("description:"):])
				}
			}

			// Check for errors during scanning
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			// Print the extracted fields
			fmt.Println("Title:", title)
			fmt.Println("Price:", price)

			// Split the text by "..."
			parts := strings.Split(description, "...")

			// Trim spaces from each part
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}

			// Join the parts with "...\n"
			formattedDescription := strings.Join(parts, "\n\n")

			fmt.Println("Description: " + formattedDescription)

			fmt.Println("past description...")

			page.MustScreenshot("process.png")

			// set title
			titleInput, err := page.Element(`label[aria-label="Title"]`)
			if err != nil {
				log.Println("Error getting title input:", err)
				return
			}

			err = titleInput.Input(title)
			if err != nil {
				log.Println("Error inputing title:", err)
				return
			}

			fmt.Println("past card...")

			// set price
			card.MustElement(`label[aria-label="Price"]`).MustInput(price)

			fmt.Println("past here...")

			// set condition
			condsInput, err := card.Element(`label[aria-label="Condition"]`)
			if err != nil {
				log.Println("Error getting conditions input:", err)
				return
			}

			err = condsInput.Click("left", 1)
			if err != nil {
				log.Println("Error clicking conditions input:", err)
				return
			}

			conds, err := page.Elements(`div[role="option"]`)
			if err != nil {
				log.Println("Error getting all conditions:", err)
				return
			}

			if len(conds) > 0 {
				err := conds[0].Click("left", 1)
				if err != nil {
					log.Println("Error selecting condition:", err)
					return
				}
			}

			fmt.Println("past here...")

			// show more details
			moreDetails, err := card.Element("div.x6s0dn4.xkh2ocl.x1q0q8m5.x1qhh985.xu3j5b3.xcfux6l.x26u7qi.xm0m39n.x13fuv20.x972fbf.x9f619.x78zum5.x1q0g3np.x1iyjqo2.xs83m0k.x1qughib.xat24cr.x11i5rnm.x1mh8g0r.xdj266r.x2lwn1j.xeuugli.x18d9i69.x4uap5.xkhd6sd.xexx8yu.x1n2onr6.x1ja2u2z")
			if err != nil {
				log.Println("Error getting more details:", err)
				return
			}

			err = moreDetails.Click("left", 1)
			if err != nil {
				log.Println("Error clicking on more details:", err)
				return
			}

			fmt.Println("past here...")

			//set description
			page.MustElementR("label", "Description").MustWaitVisible()

			textarea := page.MustElementR("label", "Description").MustElement("textarea")

			err = textarea.Input(formattedDescription)
			if err != nil {
				log.Println("Error inputing description text:", err)
				return
			}

			fmt.Println("past here...")

			// set product tags
			productTagsTextareas, err := card.Elements(`label[aria-label="Product tags"]`)
			if err != nil {
				log.Println("Error getting product tags text area:", err)
				return
			}

			// input product tags
			if len(productTagsTextareas) > 0 {
				tags := []string{
					"lipa mdogo mdogo smartphones",
					"samsung lipa mdogo mdogo",
					"iphone lipa mdogo mdogo",
					"lipa mdogo mdogo phones",
					"mkopa phones",
					"m-kopa phones",
					"lipa pole pole smartphones",
					"samsung lipa pole pole",
					"iphone lipa pole pole",
				}

				for _, tag := range tags {
					// get product tags textarea
					productTagsTextarea, err := productTagsTextareas[0].Element("textarea")
					if err != nil {
						log.Println("Error getting product tags text area:", err)
						return
					}

					// input tag
					err = productTagsTextarea.Input(tag)
					if err != nil {
						log.Println("Error inputing product tag:", err)
						return
					}
					page.Keyboard.Press(input.Enter)
				}
			}

			time.Sleep(10 * time.Second)

			// go to next page
			card.MustElement(`div[aria-label="Next"]`).MustClick()

			time.Sleep(10 * time.Second)

			// list in marketplace
			page.MustElements("div.x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.xr9ek0c.xjpr12u.xzboxd6.x14l7nz5")[2].MustClick()
			fmt.Println("listed in marketplace...")

			// list in groups
			suggestedGroups, err := page.Elements(".x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.xr9ek0c.xjpr12u.xzboxd6.x14l7nz5")
			if err != nil {
				log.Println("Error getting suggested groups wrappers:", err)
				return
			}

			groups, err := suggestedGroups[4].Elements(`div[data-visualcompletion="ignore-dynamic"]`)
			if err != nil {
				log.Println("Error getting suggested groups:", err)
				return
			}

			// calculate and log total groups
			totalGroups := len(groups)

			fmt.Printf("Total groups: %d\n", totalGroups)

			// select up to 20 groups
			if totalGroups > 20 {
				//Click on up to 20 divs randomly
				for i := 0; i < 1000 && i < len(groups); i++ {
					// Generate a random index within the range of groups slice
					randomIndex := r.Intn(len(groups))

					// Click on the div at the random index
					err := groups[randomIndex].Click("left", 1)
					if err != nil {
						log.Println("Error selecting suggested group randomly:", err)
						return
					}

					// Remove the clicked element from the slice to avoid clicking it again
					groups = append(groups[:randomIndex], groups[randomIndex+1:]...)
				}

			} else {
				// Click all groups if the total is 20 or fewer
				for i := 0; i < totalGroups; i++ {
					err := groups[i].Click("left", 1)
					if err != nil {
						log.Println("Error selecting suggested group:", err)
						return
					}
				}

			}

			// post ad
			page.MustElement(`div[aria-label="Post"]`).MustClick()

			log.Printf("%s posted successfully", title)

			time.Sleep(1 * time.Minute)
		}
		// break after one group
	}

}
