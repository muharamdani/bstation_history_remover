package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var cookies = make(map[string]string)

var headers http.Header

func main() {
	if fileExists("cookies.json") {
		loadCookiesFromFile()
		setHeaders()
		choose()
	} else {
		getCookies()
		setHeaders()
		choose()
	}
}

func loadCookiesFromFile() {
	cookieData, err := ioutil.ReadFile("cookies.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(cookieData, &cookies)
	if err != nil {
		log.Fatal(err)
	}
}

func getCookies() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("start-maximized", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var cookieNames []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.bilibili.tv/`),
		chromedp.WaitNotPresent(`.bstar-header-user__btn`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			allCookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, c := range allCookies {
				if c.Name == "bili_jct" || c.Name == "DedeUserID" || c.Name == "SESSDATA" || c.Name == "bstar-web-lang" || c.Name == "buvid3" {
					cookies[c.Name] = c.Value
					cookieNames = append(cookieNames, c.Name)
				}
			}
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	if contains(cookieNames, "bili_jct") && contains(cookieNames, "DedeUserID") && contains(cookieNames, "SESSDATA") && contains(cookieNames, "bstar-web-lang") && contains(cookieNames, "buvid3") {
		saveCookiesToFile(cookies)
	} else {
		fmt.Println("Required cookies not found")
	}
}

func setHeaders() {
	headers = http.Header{
		"Cookie": []string{
			"bili_jct=" + cookies["bili_jct"],
			"DedeUserID=" + cookies["DedeUserID"],
			"SESSDATA=" + cookies["SESSDATA"],
			"bstar-web-lang=" + cookies["bstar-web-lang"],
			"buvid3=" + cookies["buvid3"],
		},
	}
}

func saveCookiesToFile(cookies map[string]string) {
	cookieData, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("cookies.json", cookieData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func choose() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("What history do you want to delete?")
		fmt.Println("1. Anime history")
		fmt.Println("2. Video history")
		fmt.Println("3. Exit")
		fmt.Print("Your input: ")

		input, _ := reader.ReadString('\n')

		if input == "1\n" {
			deleter(4, 101)
			continue
		} else if input == "2\n" {
			deleter(3, 0)
			continue
		} else if input == "3\n" {
			break
		} else {
			fmt.Println("Error: Invalid input. Please enter 1 or 2.")
			continue
		}
	}
}

func deleter(busineesType int, subType int) {
	url := fmt.Sprintf("https://api.bilibili.tv/intl/gateway/web/v2/history/list?platform=web&business=%v&sub_type=%v&cursor=&ps=50", busineesType, subType)
	delUrl := "https://api.bilibili.tv/intl/gateway/web/v2/history/del?platform=web"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = headers
	resp, _ := client.Do(req)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()

	today := result["data"].(map[string]interface{})["today"].(map[string]interface{})["cards"].([]interface{})
	yesterday := result["data"].(map[string]interface{})["yesterday"].(map[string]interface{})["cards"].([]interface{})
	earlier := result["data"].(map[string]interface{})["earlier"].(map[string]interface{})["cards"].([]interface{})

	if len(today) == 0 && len(yesterday) == 0 && len(earlier) == 0 {
		fmt.Println("No history to delete")
		fmt.Println()
		return
	} else {
		fmt.Println("Found history to delete")
		fmt.Println()
	}

	for _, v := range today {
		card := v.(map[string]interface{})
		pld := map[string]interface{}{
			"type": busineesType,
			"oid":  card["aid"],
		}

		fmt.Println("Deleting history: ", card["title"].(string))
		pldBytes, _ := json.Marshal(pld)
		req2, _ := http.NewRequest("POST", delUrl, bytes.NewBuffer(pldBytes))
		req2.Header = headers
		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Println(err)
		} else {
			resp2.Body.Close()
		}
	}

	for _, v := range yesterday {
		card := v.(map[string]interface{})
		pld := map[string]interface{}{
			"type": busineesType,
			"oid":  card["aid"],
		}

		fmt.Println("Deleting history: ", card["title"].(string))
		pldBytes, _ := json.Marshal(pld)
		req2, _ := http.NewRequest("POST", delUrl, bytes.NewBuffer(pldBytes))
		req2.Header = headers
		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Println(err)
		} else {
			resp2.Body.Close()
		}
	}

	for _, v := range earlier {
		card := v.(map[string]interface{})
		pld := map[string]interface{}{
			"type": busineesType,
			"oid":  card["aid"],
		}

		fmt.Println("Deleting history: ", card["title"].(string))
		pldBytes, _ := json.Marshal(pld)
		req2, _ := http.NewRequest("POST", delUrl, bytes.NewBuffer(pldBytes))
		req2.Header = headers
		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Println(err)
		} else {
			resp2.Body.Close()
		}
	}
}
