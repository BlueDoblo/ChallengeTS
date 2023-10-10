package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strings"
)

// eBayItem represents the structure of each eBay item to look for information
type eBayItem struct {
	Title      string `json:"title"`
	Price      string `json:"price"`
	ProductURL string `json:"product_url"`
	Condition  string `json:"condition"`
}

var filesGenerated = 0

func main() {

	// Begin Execution Parameters
	// Allowed values are
	// 0 : Process all the items
	// 0 : Process all the items
	// 1 : Only crawl new items with condition
	// 2 : Only crawl new pre-owned items
	var parameterIndex = "0"                // allowed values 0 or 1 or 2
	var parameterValue = "Totalmente nuevo" // allowed values "Totalmente nuevo" or "De segunda mano"
	var flagContinue = false

	// End Execution Parameters

	// Creation of the collector of the information
	c := colly.NewCollector()

	// Set up a folder to store results
	outputFolder := "data"

	// Create the output folder if it doesn't exist
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {

		os.Mkdir(outputFolder, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating the folder:", err)
		}

	}

	// Define a callback to be executed when a list item is found
	c.OnHTML(".s-item", func(e *colly.HTMLElement) {
		// Extract the title, price, and product URL
		title := e.ChildText(".s-item__title")
		price := e.ChildText(".s-item__price")
		productURL := e.ChildAttr("a.s-item__link", "href")
		condition := e.ChildText(".s-item__subtitle")

		// Extract the item ID from the product URL
		itemID := extractItemID(productURL)

		// Create a struct with the extracted data
		item := eBayItem{
			Title:      title,
			Price:      price,
			ProductURL: productURL,
			Condition:  condition,
		}

		// Convert the item to JSON
		jsonData, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		// Create a file with the item ID as the filename in the output folder
		index := strings.Index(itemID, "?")

		// Evualuate the second alternative
		if strings.Compare(parameterIndex, "0") == 0 {
			flagContinue = true
		}

		// Evualuate the first alternative
		if strings.Compare(parameterIndex, "1") == 0 &&
			strings.Compare(condition, parameterValue) == 0 {
			flagContinue = true
		} else {
			if strings.Compare(parameterIndex, "1") == 0 {
				flagContinue = false
			}
		}

		// Evualuate the second alternative
		if strings.Compare(parameterIndex, "2") == 0 &&
			strings.Compare(condition, parameterValue) == 0 {
			flagContinue = true
		} else {
			if strings.Compare(parameterIndex, "2") == 0 {
				flagContinue = false
			}
		}

		// Verify the existence of an item id
		if index != -1 && flagContinue {

			// If i found the symbol ? i need to format the name to use for the file name
			result := itemID[:index]
			filename := fmt.Sprintf("%s/%s.json", outputFolder, result)
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			// Write the JSON data to the file
			_, err = file.Write(jsonData)
			if err != nil {
				fmt.Println("Error closing the file:", err)
				return
			}
			filesGenerated++
			//fmt.Printf("Saved data for item ID %s\n", result)
		}
	})

	// Start the crawling process from the eBay page
	url := "https://www.ebay.com/sch/garlandcomputer/m.html"
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error visiting eBay page:", err)
	}

	fmt.Println("Files Generated: ", filesGenerated)
}

// extractItemID extracts the item ID from a product URL
func extractItemID(productURL string) string {
	parts := strings.Split(productURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
