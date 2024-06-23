package main

import (
	"stressTestDNS/dnsquerier"
	"stressTestDNS/domainfetcher"
)

func main() {
	domainfetcher.FetchAndSaveDomains()
	dnsquerier.LoadAndQueryDomains()
}
