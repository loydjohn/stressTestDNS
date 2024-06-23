package dnsquerier

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const (
	Rate            = 5
	DomainFile      = "domains.txt" // Change this to your domain list file path
	TargetDNSServer = `8.8.8.8:53`
)

// Load valid domains from a file
func LoadValidDomains(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return domains, nil
}

// Perform DNS query
func queryDNS(domain string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	_, _, err := client.Exchange(message, TargetDNSServer)

	if err != nil {
		results <- fmt.Sprintf("%s: error - %v", domain, err)
	} else {
		results <- fmt.Sprintf("%s: resolved", domain)
	}
}

// Function to manage DNS queries
func QueryManager(domains []string) {
	var wg sync.WaitGroup
	results := make(chan string, Rate)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	for range ticker.C {
		rand.Shuffle(len(domains), func(i, j int) { domains[i], domains[j] = domains[j], domains[i] })

		for i := 0; i < Rate && i < len(domains); i++ {
			wg.Add(1)
			go queryDNS(domains[i], &wg, results)
		}
		wg.Wait()
	}
}

func LoadAndQueryDomains() {
	domains, err := LoadValidDomains(DomainFile)
	if err != nil {
		log.Fatalf("Failed to load valid domains: %v", err)
	}

	QueryManager(domains)
}
