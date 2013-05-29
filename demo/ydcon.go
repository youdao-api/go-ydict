package main

import (
	"bufio"
	"fmt"
	"github.com/daviddengcn/go-ydict"
	"os"
	"strings"
)

func wordChars(t string) int {
	cnt := 0
	for _, c := range t {
		cnt++
		if c >= 128 {
			cnt++
		}
	}
	return cnt
}

func main() {
	client := ydict.NewOnlineClient("go-ydict", "252639882")

	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Query: ")
		q, err := in.ReadString('\n')
		if err != nil {
			break
		}
		res, err := client.Query(q)
		if err != nil {
			fmt.Println("Query failed:", err)
			continue
		}

		fmt.Println(res.Query)
		fmt.Println(strings.Repeat("-", wordChars(res.Query)))
		if res.Translation != nil {
			fmt.Println("翻译")
			for _, t := range res.Translation {
				fmt.Printf("    %s\n", t)
			}
		}

		if res.Basic != nil {
			fmt.Println("基本释义")
			if res.Basic.Phonetic != "" {
				fmt.Printf("    /%s/\n", res.Basic.Phonetic)
			}
			for _, exp := range res.Basic.Explains {
				fmt.Printf("    %s\n", exp)
			}
		}

		if len(res.Web) > 0 {
			fmt.Println("网络释义")
			wd := 0
			for _, web := range res.Web {
				wc := wordChars(web.Key)
				if wc > wd {
					wd = wc
				}
			}
			for _, web := range res.Web {
				fmt.Printf("    %s%*s  %s\n", web.Key, wd-wordChars(web.Key),
					"", strings.Join(web.Value, "; "))
			}
		}
	}
}
