package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"suzuki/webscraping"
)

func Hashstr(Txt string) string {
	h := sha1.New()
	h.Write([]byte(Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x\n", bs))
	return sh
}

func main() {
	//"MA3FB21S080969810"
	var vin string
	flag.StringVar(&vin, "vin", "MA3FB21S080969810", "VIN A BUSCAR")
	flag.Parse()
	suzukiBot := webscraping.BotSuzuki{}
	suzukiBot.User = "MMoron01"
	suzukiBot.Pass = "MMoron01"
	suzukiBot.Init(vin)

}
