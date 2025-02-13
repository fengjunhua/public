package main

import (
	"fmt"
	"github.com/8zhiniao/public/utils"
)

func main() {

	name, err := utils.GetUserName()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(name)
}
