package config

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"errors"
)

type Config struct  {
	DbUrl	 string `json:"db_url"`
	CurrUserName	 string `json:"current_user_name"`
}

type State struct {
	ConfPtr *Config
}

type Command struct {
	CommandName string
	Args [] string
}

type Commands struct {
	Handlers map[string]func(*State,Command) error //map of command names to their handler functiosn
}

const configFileName = "/.gatorconfig.json"

func (c *Commands) Register(name string, f func(*State, Command) error) {
	//registers a new handler function for a command name
	c.Handlers[name] = f
	return 
}


func (c *Commands) Run(s *State, cmd Command) error {
	//runs a given command with the provided state if it exists
	cmdHandler, ok := c.Handlers[cmd.CommandName]
	if !ok {
		return errors.New("error was not found in commands map")
	}
	return cmdHandler(s,cmd)
}

func HandlerLogin(s *State, cmd Command) error {
	fmt.Println("in handlerlogin")
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("username required for login command") 
	}
	err := s.ConfPtr.SetUser(cmd.Args[0])
	if err != nil {
		return err 
	}
	fmt.Println("user has been set to %s\n", cmd.Args[0])
	return nil
}



func getConfigFilePath() (string, error) {

	homeDir,err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error with getting home dir")
		return "",err
	}
	fullDir := filepath.Join(homeDir + configFileName)
	return fullDir,nil 

}
func Read() (Config, error) {
//read home DIR, then decode to JSON string into a new Config struct
	
	var conf Config

	fullDir,err := getConfigFilePath()
	if err != nil {
		fmt.Println("error: Couldn't get config file path",err)
	}

	jsonFile, err := os.Open(fullDir)

	if err != nil {
		fmt.Println(err)
		return conf, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		fmt.Println("error:",err)
		return conf,err
	}
	return conf,nil

}

func (conf *Config) SetUser(user string) error {
	

	conf.CurrUserName = user
	fmt.Println("user is now ",conf.CurrUserName)
	return write(*conf)
}

func write(conf Config) error {



	fmt.Println("conf is ",conf)
	jsonData, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	fmt.Println("jsondata is ",string(jsonData))
	fullDir,err := getConfigFilePath()
	if err != nil {
		return err
	
	}
	fmt.Println("fullDir is ",fullDir)
	err = os.WriteFile(fullDir,jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
