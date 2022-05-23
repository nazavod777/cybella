package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var client *fasthttp.Client

func main() {
	clear()
	fmt.Print("Telegram channel - t.me/n4z4v0d\n\n")
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	var emailFolder string
	var useProxies string

	fmt.Print("Drop .txt with emails: ")
	fmt.Scanf("%s\n", &emailFolder)

	fmt.Print("Use Tor Proxies? (y/N): ")
	fmt.Scanf("%s\n", &useProxies)

	file, err := os.Open(emailFolder)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	clear()

	for scanner.Scan() {
		wg.Add(1)
		data := scanner.Text()
		go func(data string, useProxies string) {
			sendRegisterReq(data, useProxies)
			defer wg.Done()
		}(data, useProxies)
	}
	wg.Wait()
	fmt.Printf("\nThe work was successfully completed\n")
	fmt.Println("Press Enter To Exit...")
	var input string
	fmt.Scanf("%s", &input)
	os.Exit(0)
}

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func RandomString(n int) string {
	var letters = []rune("0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func writeResult(email string, status bool) {
	if status {
		file, err := os.OpenFile("good.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			file.WriteString(email + "\n")
		} else {
			fmt.Println(err)
		}
	} else if !status {
		file, err := os.OpenFile("error.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			file.WriteString(email + "\n")
		} else {
			fmt.Println(err)
		}
	}
	return
}

func sendRegisterReq(email string, useProxies string) {
	for i := 0; i < 11; i++ {
		jQuery_rand := RandomString(20)
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)

		readTimeout, _ := time.ParseDuration("5s")
		writeTimeout, _ := time.ParseDuration("5s")
		maxIdleConnDuration, _ := time.ParseDuration("1h")

		if useProxies == "y" {
			tor_user := rand.Int()
			tor_pass := rand.Int()

			client = &fasthttp.Client{
				ReadTimeout:                   readTimeout,
				WriteTimeout:                  writeTimeout,
				MaxIdleConnDuration:           maxIdleConnDuration,
				NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
				DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
				DisablePathNormalizing:        true,
				// increase DNS cache time to an hour instead of default minute
				Dial: fasthttpproxy.FasthttpSocksDialer("socks5://" + strconv.Itoa(tor_user) + ":" + strconv.Itoa(tor_pass) + "@localhost:9150"),
			}
		} else {
			client = &fasthttp.Client{
				ReadTimeout:                   readTimeout,
				WriteTimeout:                  writeTimeout,
				MaxIdleConnDuration:           maxIdleConnDuration,
				NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
				DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
				DisablePathNormalizing:        true,
			}
		}

		req := fasthttp.AcquireRequest()
		req.Header.Set("accept", "*/*")
		req.Header.Set("referer", "https://www.cybella.xyz/")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7")
		req.SetRequestURI("https://gmail.us20.list-manage.com/subscribe/post-json?u=4ba44547646dfe80e06e149cc&amp;id=4d3f338d6d&c=jQuery" + jQuery_rand + "_" + timestamp + "&email=" + email + "&EMAIL=" + email + "&b_4ba44547646dfe80e06e149cc_4d3f338d6d=&_=" + timestamp)
		req.Header.SetMethod(fasthttp.MethodGet)
		resp := fasthttp.AcquireResponse()
		err := client.Do(req, resp)

		if err == nil && strings.Contains(fmt.Sprintf("%s", resp.Body()), `result":"success",`) {
			fmt.Printf("[%s] | The account has been successfully registered: %s\n", email, resp.Body())
			writeResult(email, true)
			fasthttp.ReleaseResponse(resp)
			return
		} else if err == nil && strings.Contains(fmt.Sprintf("%s", resp.Body()), `result":"error",`) {
			fmt.Printf("[%s] | Wrong Response: %s\n", email, resp.Body())
		} else {
			fmt.Fprintf(os.Stderr, "[%s] | Unexpected error: %v\n", email, err)
		}
		fasthttp.ReleaseResponse(resp)
	}
	writeResult(email, false)
	return
}
