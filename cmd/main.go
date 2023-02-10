package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cluttrdev/deepl-go/pkg/deepl"
)

func main() {
	timeout := 10 * time.Second
	client := deepl.NewClient(deepl.BaseURLFree, os.Getenv("DEEPL_API_KEY"), timeout)

	translations, err := client.TranslateText([]string{"Hello, world!"}, "DE")
	if err != nil {
		log.Fatal(err)
	}

	for _, translation := range translations {
		fmt.Println(translation.Text)
	}
}
