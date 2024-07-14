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

// 	time.Sleep(1 * time.Minute)

// 	page = browser.MustPage("https://web.facebook.com/").MustWaitLoad()

// 	page.MustScreenshot("facebook.png")
// }

package main

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
)

func main() {
	// Launch the browser with specific configurations
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // Ensure the browser closes when main function exits

	// Root directory containing subdirectories with images
	root := "/home/kwandapchumba/Pictures/MO-PHONES"
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading root directory:", err)
		return
	}

	// Loop through each subdirectory in the root directory
	for _, entry := range entries {
		if entry.IsDir() {
			// Open the Facebook Marketplace item creation page
			page := browser.MustPage("https://web.facebook.com/marketplace/create/item").MustWaitLoad()

			// Locate the file input element on the page
			fileInput := page.MustElement(`input[type="file"]`)

			// Path to the current subdirectory
			subDir := filepath.Join(root, entry.Name())
			fmt.Println("Directory:", subDir)

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
						fmt.Println("File:", filePath)

						// Collect file paths to attach
						imageFiles = append(imageFiles, filePath)
					}
				}
			}

			// Attach collected files to the file input element
			if len(imageFiles) > 0 {
				log.Println(imageFiles)
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
				if strings.HasPrefix(line, "title:") {
					title = strings.TrimSpace(line[len("title:"):])
				} else if strings.HasPrefix(line, "price:") {
					price = strings.TrimSpace(line[len("price:"):])
				} else if strings.HasPrefix(line, "description:") {
					description = strings.TrimSpace(line[len("description:"):])
				}
			}

			// Check for errors during scanning
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
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
			formattedDescription := strings.Join(parts, ".\n\n")

			fmt.Println("Description: " + formattedDescription)

			// get inputs wrapper
			marketplace := page.MustElement(`div[aria-label="Marketplace"]`)

			// title
			t := marketplace.MustElement(`label[aria-label="Title"]`)
			t.MustElement("input.x1i10hfl").MustInput(title)

			// price
			p := marketplace.MustElement(`label[aria-label="Price"]`)
			p.MustElement("input.x1i10hfl").MustInput(price)

			// category
			marketplace.MustElement(`label[aria-label="Category"]`).MustClick()

			// Find all parent divs with data-visualcompletion="ignore-dynamic"
			categoryParentDivs := page.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

			if len(categoryParentDivs) > 0 {
				categoryParentDivs[25].MustClick()

			}

			// condition
			marketplace.MustElement(`label[aria-label="Condition"]`).MustClick()
			conditionParentDivs := page.MustElements(`div[role="option"]`)
			if len(conditionParentDivs) > 0 {
				conditionParentDivs[0].MustClick()

			}

			// desc
			// Find all div elements that match the specified structure and classes
			divs := page.MustElements(`div.x1pi30zi.x1swvt13.xyamay9`)
			if len(divs) > 0 {
				divs[3].MustElement("textarea").MustInput(formattedDescription).MustFocus()
			}

			// availability
			// Define the label selector
			availabilitySelector := `label[aria-label="Availability"]`

			// Find all label elements with the specified structure
			labels := page.MustElements(availabilitySelector)

			if len(labels) > 0 {
				// for _, label := range labels {

				// }
				labels[0].MustClick()
			}

			// Define the div selector
			divSelector := `div.html-div.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x6s0dn4.x78zum5.x1q0g3np.x1iyjqo2.x1qughib.xeuugli`

			// Find all div elements with the specified structure
			availabilityDivs := page.MustElements(divSelector)

			if len(availabilityDivs) > 0 {
				availabilityDivs[1].MustClick()
			}

			// Product tags
			productTagDivs := page.MustElements(`label[aria-label="Product tags"]`)
			if len(productTagDivs) > 0 {
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
				// productTagDivs[0].MustElement("textarea").MustInput("lipa mdogo mdogo phones")
				for _, t := range tags {
					productTagDivs[0].MustElement("textarea").MustInput(t)
					page.Keyboard.Press(input.Enter)
				}

			}

			// Location
			// location := page.MustElement(`label[aria-label="Location"]`)
			// location.MustClick()

			//time.Sleep(10 * time.Second)

			// next button
			page.MustElement(`div[aria-label="Next"]`).MustClick()

			time.Sleep(10 * time.Second)

			// groups
			wrapper := page.MustElements(".x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.xr9ek0c.xjpr12u.xzboxd6.x14l7nz5")

			groupsContainer := wrapper[3]

			groups := groupsContainer.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			totalGroups := len(groups)

			fmt.Printf("Total groups: %d\n", totalGroups)

			if totalGroups > 20 {
				//Click on up to 20 divs randomly
				for i := 0; i < 5000 && i < len(groups); i++ {
					// Generate a random index within the range of groups slice
					randomIndex := r.Intn(len(groups))

					// Click on the div at the random index
					groups[randomIndex].MustClick()

					// Remove the clicked element from the slice to avoid clicking it again
					groups = append(groups[:randomIndex], groups[randomIndex+1:]...)
				}

			} else {
				// Click all groups if the total is 20 or fewer
				for i := 0; i < totalGroups; i++ {
					groups[i].MustClick()
				}

			}

			// publish
			page.MustElement(`div[aria-label="Publish"]`).MustClick()

			time.Sleep(30 * time.Second)

			page.MustScreenshotFullPage("facebook.png")
		}
	}
}
