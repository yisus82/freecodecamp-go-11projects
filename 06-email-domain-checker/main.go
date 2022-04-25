package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	fmt.Print("Enter a domain: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		if domain == "" {
			break
		}
		checkDomain(domain)
		fmt.Print("\nEnter a domain: ")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mx, err := net.LookupMX(domain)
	if err != nil {
		color.Red("Error: %v\n", err)
	}
	if len(mx) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		color.Red("Error: %v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		color.Red("Error: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	color.Green("Domain: %s\n", domain)
	color.Green("Has MX?: %t\n", hasMX)
	color.Green("Has SPF?: %t\n", hasSPF)
	color.Green("SPF record: %v\n", spfRecord)
	color.Green("Has DMARC?: %t\n", hasDMARC)
	color.Green("DMARC record: %v\n", dmarcRecord)
}
