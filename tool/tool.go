package main

import (
	"fmt"

	"github.com/tbuckley/whistler"
)

func main() {
	whistler.Initialize()
	defer whistler.Terminate()

	whistle, err := whistler.New()
	if err != nil {
		panic(err)
	}
	defer whistle.Close()

	kikeeChan := whistle.Add(whistler.Kikee)
	go func() {
		for {
			<-kikeeChan
			fmt.Println("Whistle whistle!")
		}
	}()

	whistle.Listen()

	fmt.Println("Press ENTER to quit")
	_, err = fmt.Scanln()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}
}
