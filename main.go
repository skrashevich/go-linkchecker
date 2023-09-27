package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/samber/lo"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	files       string
	allowErrors string
	// allowDupe     bool
	allowSSL      bool
	allowRedirect bool
	allowTimeout  bool
	baseURL       string
	requestDelay  int
	timeout       int
	skipSave      bool
	whiteList     string
)

type checkResult struct {
	Link string
	OK   bool
}

func checkLink(link string, ch chan checkResult, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: allowSSL},
		DisableKeepAlives: true,
	}
	var checkRedirect func(req *http.Request, via []*http.Request) error
	if allowRedirect {
		checkRedirect = nil
	} else {
		checkRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	client := &http.Client{
		Transport:     tr,
		Timeout:       time.Duration(timeout) * time.Second,
		CheckRedirect: checkRedirect,
	}
	resp, err := client.Get(link)

	if err != nil || resp.StatusCode != 200 {
		ch <- checkResult{Link: link, OK: false}
		return

	}
	ch <- checkResult{Link: link, OK: true}
}

func createWhiteListSet(whiteList string) map[string]bool {
	whiteListSet := make(map[string]bool)
	if whiteList == "" {
		return whiteListSet
	}
	for _, url := range strings.Split(strings.ReplaceAll(whiteList, " ", ""), ",") {
		whiteListSet[url] = true
	}
	return whiteListSet
}

func main() {

	flag.StringVar(&files, "f", "", "Comma separated files to check")
	flag.StringVar(&files, "files", "", "Comma separated files to check")
	flag.StringVar(&allowErrors, "a", "", "Status code errors to allow")
	flag.StringVar(&allowErrors, "allow", "", "Status code errors to allow")
	// flag.BoolVar(&allowDupe, "allow-dupe", false, "Duplicate URLs are allowed")
	flag.BoolVar(&allowSSL, "allow-ssl", false, "SSL errors are allowed")
	flag.BoolVar(&allowRedirect, "allow-redirect", false, "Redirected URLs are allowed")
	flag.BoolVar(&allowTimeout, "allow-timeout", false, "URLs that time out are allowed")
	flag.StringVar(&baseURL, "base-url", "", "Base URL to use for relative links")
	flag.IntVar(&requestDelay, "d", 0, "Set request delay")
	flag.IntVar(&requestDelay, "request-delay", 0, "Set request delay")
	flag.IntVar(&timeout, "t", 30, "Set connection timeout")
	flag.IntVar(&timeout, "set-timeout", 30, "Set connection timeout")
	flag.BoolVar(&skipSave, "skip-save-results", false, "Skip saving results")
	flag.StringVar(&whiteList, "w", "", "Comma separated URLs to white list")
	flag.StringVar(&whiteList, "white-list", "", "Comma separated URLs to white list")
	flag.Parse()

	args := flag.Args() // Non-flag arguments

	startTime := time.Now()
	fileList := strings.Split(files, ",")
	fileList = append(fileList, args...)

	whiteListSet := createWhiteListSet(whiteList)

	linkRegex := regexp.MustCompile(`\((http[^\)]+)\)`)
	ch := make(chan checkResult)
	var wg sync.WaitGroup

	for _, filePath := range fileList {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		alllinks := make([]string, 0)
		totalLinks := 0
		for scanner.Scan() {
			line := scanner.Text()
			matches := linkRegex.FindAllStringSubmatch(line, -1)

			for _, match := range matches {

				alllinks = append(alllinks, match[1])

				totalLinks++
			}
		}
		links := lo.Uniq(alllinks)
		fmt.Printf("Links to check: %d, %d unique\n\n", totalLinks, len(links))
		i := 0
		for _, link := range links {
			i++
			fmt.Printf("\t%d. %s\n", i, link)
		}
		for _, link := range links {
			if whiteListSet[link] {
				continue
			}
			wg.Add(1)
			go checkLink(link, ch, &wg)
			if requestDelay > 0 {
				time.Sleep(time.Duration(requestDelay) * time.Second)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
	}

	fmt.Print("Checking URLs: ")

	go func() {
		wg.Wait()
		close(ch)
	}()

	brokenlist := make([]string, 0)
	for msg := range ch {
		if !msg.OK {
			fmt.Print("x")
			brokenlist = append(brokenlist, msg.Link)
		}
		fmt.Print("✓")
	}
	fmt.Println("\nBroken links:")
	fmt.Println()
	for _, link := range brokenlist {
		fmt.Println(link)
	}
	fmt.Println()

	elapsedTime := time.Since(startTime)

	fmt.Printf("Seconds elapsed: %v\n", elapsedTime)

}
