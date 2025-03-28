package main

import (
	"fmt"
	"github.com/willmelton21/gator/internal/config"
	"os"
)

func main() {

	conf,err := config.Read()
	
	if err != nil {
		fmt.Println("err from read is :",err)
	}

	state := &config.State{ConfPtr: &conf}
	cmds := &config.Commands{
		Handlers: make(map[string]func(*config.State, config.Command) error),
	}
        
	cmds.Register("login",config.HandlerLogin)

	args := os.Args
	if (len(args) < 2) {
		fmt.Println("Error: Not enough arguments. Usage: gator <command> [args]")		
		os.Exit(1)
	}

	commandName := args[1]

	var commandArgs []string
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	cmd := config.Command {
		CommandName: commandName,
		Args: commandArgs,
	}
   	fmt.Println("args are ",args[2:])	
	err = cmds.Run(state,cmd)
	if err != nil {

		fmt.Println("Error:",err)
		os.Exit(1)
	}
}
