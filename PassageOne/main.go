package main

import (
	"PassageOne/cmd"
	"fmt"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err : %v", err)
	}
	fmt.Println("谢谢使用")
}
