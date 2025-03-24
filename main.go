package main

import (
	"fmt"
	"github.com/willmelton21/gator/internal/config"
)

func main() {

	conf,err := config.Read()

	if err != nil {
		fmt.Println("err from read is :",err)
	}

	conf.SetUser("Will")
	newConf,err := config.Read()
	if err != nil {
		fmt.Println("err from read is :",err)
	}

	fmt.Printf("%+v\n", newConf)

}
