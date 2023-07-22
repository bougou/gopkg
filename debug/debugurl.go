package debug

import (
	"fmt"
	"net/url"
)

func DebugURL(urlstr string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("url scheme:", u.Scheme)
	fmt.Println("url domain:", u.Host)
	fmt.Println("url domain w/o port:", u.Hostname())
	fmt.Println("url path:", u.Path)
	fmt.Println("url raw path:", u.RawPath)
}
