package config

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Config struct  {
	DbUrl	 string `json:"db_url"`
	currUserName	 string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

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
	

	conf.currUserName = user
	fmt.Println("user is now ",conf.currUserName)
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
