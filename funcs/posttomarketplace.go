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

func PostToMarketplace() {
	// load config files
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	// Launch the browser with specific configurations
	u := launcher.New().UserDataDir("~/.config/google-chrome").Leakless(true).NoSandbox(true).Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // Ensure the browser closes when main function exits

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

	for _, entry := range entries {
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
			fmt.Println(imageFiles)
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
		formattedDescription := strings.Join(parts, ".\n\n")

		fmt.Println("Description: " + formattedDescription)

		// get inputs wrapper
		marketplace, err := page.Element(`div[aria-label="Marketplace"]`)
		if err != nil {
			log.Println("Error getting marketplace div:", err)
			return
		}

		// get title input
		titleInput, err := marketplace.Element(`label[aria-label="Title"]`)
		if err != nil {
			log.Println("Error getting title input:", err)
			return
		}

		// input title text
		err = titleInput.Input(title)
		if err != nil {
			log.Println("Error inputing title text:", err)
			return
		}

		// get price input
		priceInput, err := marketplace.Element(`label[aria-label="Price"]`)
		if err != nil {
			log.Println("Error getting price input:", err)
			return
		}

		// input price text
		err = priceInput.Input(price)
		if err != nil {
			log.Println("Error inputing price text:", err)
			return
		}

		// get category input
		categoryInput, err := marketplace.Element(`label[aria-label="Category"]`)
		if err != nil {
			log.Println("Error getting category input:", err)
			return
		}

		// click category input
		err = categoryInput.Click("left", 1)
		if err != nil {
			log.Println("Error clicking category input:", err)
			return
		}

		// Find all parent divs with data-visualcompletion="ignore-dynamic"
		categoryParentDivs := page.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

		if len(categoryParentDivs) > 0 {
			categoryParentDivs[25].MustClick()

		}

		// get condition input
		conditionInput, err := marketplace.Element(`label[aria-label="Condition"]`)
		if err != nil {
			log.Println("Error getting condition input:", err)
			return
		}

		// click conditionInput
		err = conditionInput.Click("left", 1)
		if err != nil {
			log.Println("Error clicking condition input:", err)
			return
		}

		// get all conditions
		conditions, err := page.Elements(`div[role="option"]`)
		if err != nil {
			log.Println("Error getting all conditions:", err)
			return
		}

		// select neccessary condition
		if len(conditions) > 0 {
			err := conditions[0].Click("left", 1)
			if err != nil {
				log.Println("Error selecting condition:", err)
				return
			}
		}

		// get description label
		page.MustElementR("label", "Description").MustWaitVisible()

		// Locate the textarea within the label
		textarea := page.MustElementR("label", "Description").MustElement("textarea")

		// input description text
		err = textarea.Input(formattedDescription)
		if err != nil {
			log.Println("Error inputing description text:", err)
			return
		}

		// get availability inputs
		availabilitySelector := `label[aria-label="Availability"]`
		availabilityInputs, err := page.Elements(availabilitySelector)
		if err != nil {
			log.Println("Error getting availability inputs:", err)
			return
		}

		// open available options
		if len(availabilityInputs) > 0 {
			err := availabilityInputs[0].Click("left", 1)
			if err != nil {
				log.Println("Error opening available options:", err)
				return
			}
		}

		// find all available options
		divSelector := `div.html-div.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x6s0dn4.x78zum5.x1q0g3np.x1iyjqo2.x1qughib.xeuugli`

		availableOptions, err := page.Elements(divSelector)
		if err != nil {
			log.Println("Error getting available options:", err)
			return
		}

		// click neccessary availability option
		if len(availableOptions) > 0 {
			err := availableOptions[1].Click("left", 1)
			if err != nil {
				log.Println("Error selecting available option:", err)
				return
			}
		}

		// get product tags textareas
		productTagsTextareas, err := page.Elements(`label[aria-label="Product tags"]`)
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

		// get "next" button
		nextButton, err := page.Element(`div[aria-label="Next"]`)
		if err != nil {
			log.Println("Error getting next button:", err)
			return
		}

		// click "next" button
		err = nextButton.Click("left", 1)
		if err != nil {
			log.Println("Error clicking next button:", err)
			return
		}

		// wait 10 seconds for the next page to load
		time.Sleep(10 * time.Second)

		// get suggestedGroupsWrapper
		suggestedGroupsWrappers, err := page.Elements(".x9f619.x1n2onr6.x1ja2u2z.x78zum5.xdt5ytf.x2lah0s.x193iq5w.xr9ek0c.xjpr12u.xzboxd6.x14l7nz5")
		if err != nil {
			log.Println("Error getting suggested groups wrappers:", err)
			return
		}

		suggestedGroupsWrapper := suggestedGroupsWrappers[3]

		// get suggested groups
		groups, err := suggestedGroupsWrapper.Elements(`div[data-visualcompletion="ignore-dynamic"]`)
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

		// get publish button
		publishButton, err := page.Element(`div[aria-label="Publish"]`)
		if err != nil {
			log.Println("Error getting publish button:", err)
			return
		}

		// click publish button to publish ad
		err = publishButton.Click("left", 1)
		if err != nil {
			log.Println("Error publishing ad:", err)
			return
		}
	}
}
