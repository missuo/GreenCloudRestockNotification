/*
 * @Author: Vincent Young
 * @Date: 2023-10-16 03:16:09
 * @LastEditors: Vincent Young
 * @LastEditTime: 2023-10-18 10:34:42
 * @FilePath: /GreenCloudRestockNotification/main.go
 * @Telegram: https://t.me/missuo
 *
 * Copyright © 2023 by Vincent, All Rights Reserved.
 */
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"flag"

	"github.com/gocolly/colly"
)

func sendToBark(baseURL string, message string) {
	title := "GreenCloud is restocked, buy it now!"
	endpoint := fmt.Sprintf("%s/%s/%s", baseURL, title, message)
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received HTTP %d from Bark API\n", resp.StatusCode)
	} else {
		fmt.Println("Message sent successfully!")
	}
}

func trackGreenCloud(baseURL string) {
	var stock100usd string

	c := colly.NewCollector(
		colly.AllowedDomains("greencloudvps.com"),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	})

	c.OnHTML("#product75 > strong:nth-child(1) > header > span.qty", func(e *colly.HTMLElement) {
		stock100usd = e.Text
		stock100usd = strings.Replace(stock100usd, "Available", "", 1)
		stock100usd = strings.Join(strings.Fields(stock100usd), "")
		fmt.Println("1010 Birthday JP Stock: " + stock100usd)
		if stock100usd != "0" {
			sendToBark(baseURL, "1010 Birthday JP Stock: " + stock100usd)
		}
	})

	c.Visit("https://greencloudvps.com/billing/store/10th-birthday-sale")
}

func main() {
	// Define flags
	baseURL := flag.String("u", "", "Base URL for the Bark service")
	sleepTime := flag.Int("t", 3, "Sleep time in seconds between checks")

	// Parse the flags
	flag.Parse()

	// Check if baseURL is not empty
	if *baseURL == "" {
		fmt.Println("Error: baseURL is required. Use -u flag to set the baseURL.")
		return
	}

	sleepDuration := time.Duration(*sleepTime) * time.Second

	for {
		trackGreenCloud(*baseURL)
		time.Sleep(sleepDuration)
	}
}
