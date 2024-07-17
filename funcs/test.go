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

func Test() {
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
		var title, price, category, condition, description string

		// Create a new scanner to read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "title:"):
				title = strings.TrimSpace(line[len("title:"):])
			case strings.HasPrefix(line, "price:"):
				price = strings.TrimSpace(line[len("price:"):])
			case strings.HasPrefix(line, "category"):
				category = strings.TrimSpace(line[len("category:"):])
			case strings.HasPrefix(line, "condition"):
				condition = strings.TrimSpace(line[len("condition:"):])
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
		fmt.Println("Category:", category)
		fmt.Println("Condition:", condition)

		// Split the text by "..."
		parts := strings.Split(description, "...")

		// Trim spaces from each part
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		// Join the parts with "...\n"
		formattedDescription := strings.Join(parts, "\n\n")

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
		cats := page.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

		for _, cat := range cats {
			if cat.MustText() == category {
				cat.MustClick()
			}
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

		for _, cond := range conditions {
			if cond.MustText() == condition {
				cond.MustClick()
			}
		}

		time.Sleep(10 * time.Second)
		page.MustScreenshot("facebook.png")
	}

}
