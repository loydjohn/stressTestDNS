package domainfetcher

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	OutputFile = "domains.txt"
	URL        = "https://downloads.majestic.com/majestic_million.csv"
)

// Fetch the top 1000 domains from the Majestic Million list
func FetchTopDomains(url string) ([]string, error) {
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var domains []string
	scanner := bufio.NewScanner(strings.NewReader(resp.String()))
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		if count == 0 {
			// Skip the header line
			count++
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) > 2 {
			domains = append(domains, fields[2])
		}
		count++
		if count > 10000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// Save domains to a file
func SaveDomainsToFile(domains []string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, domain := range domains {
		_, err := writer.WriteString(domain + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func FetchAndSaveDomains() {
	domains, err := FetchTopDomains(URL)
	if err != nil {
		log.Fatalf("Failed to fetch top domains: %v", err)
	}

	err = SaveDomainsToFile(domains, OutputFile)
	if err != nil {
		log.Fatalf("Failed to save domains to file: %v", err)
	}

	fmt.Printf("Successfully saved %d domains to %s\n", len(domains), OutputFile)
}
