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
	"github.com/haron1996/fb/config"
)

func Post2Groups(page *rod.Page) {
	fmt.Println("post to group using discussion tab")
	class := `div.x1i10hfl.x1ejq31n.xd10rxx.x1sy0etr.x17r0tee.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x16tdsg8.x1hl2dhg.xggy1nq.x87ps6o.x1lku1pv.x1a2a7pz.x6s0dn4.xmjcpbm.x107yiy2.xv8uw2v.x1tfwpuw.x2g32xy.x78zum5.x1q0g3np.x1iyjqo2.x1nhvcw1.x1n2onr6.xt7dq6l.x1ba4aug.x1y1aw1k.xn6708d.xwib8y2.x1ye3gou`
	time.Sleep(10 * time.Second)
	pageHasWriteSomethingBtn, writeSomethingBtn, err := page.Has(class)
	if err != nil {
		log.Println("Error checking if page has write something button:", err)
		return
	}

	if !pageHasWriteSomethingBtn {
		return
	}

	err = writeSomethingBtn.Click("left", 1)
	if err != nil {
		log.Println("Error clicking write something button:", err)
		return
	}

	time.Sleep(10 * time.Second)

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

	for _, entry := range entries {
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

		fileInput, err := page.Element(`input[type="file"]`)
		if err != nil {
			log.Println("Error getting file input:", err)
			return
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

		time.Sleep(10 * time.Second)

		// pageHasCreatePostInputWrapper, createPostInputWrapper, err := page.Has(`div.xb57i2i.x1q594ok.x5lxg6s.x6ikm8r.x1ja2u2z.x1pq812k.x1rohswg.xfk6m8.x1yqm8si.xjx87ck.xx8ngbg.xwo3gff.x1n2onr6.x1oyok0e.x1odjw0f.x1e4zzel.x78zum5.xdt5ytf.x1iyjqo2`)
		// if err != nil {
		// 	log.Println("Error checking if page has create post input wrapper:", err)
		// 	return
		// }

		// if !pageHasCreatePostInputWrapper {
		// 	log.Println("page has no create post input wrapper")
		// 	return
		// }

		// time.Sleep(10 * time.Second)

		// log.Println(createPostInputWrapper.Has(`div[role="textbox"]`))
		// Find the element by its aria-label
		log.Println(page.Has(`div[aria-describedby="placeholder-8mo26"]`))
	}

	time.Sleep(10 * time.Second)
	page.MustScreenshot("facebook.png")
}
