package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	DbUserName string `json:"db_user_name"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbIP       string `json:"db_ip"`
}

func Parse() Config {
	JsonConfig, err := os.Open(".\\config\\config.json")
	config := Config{}
	//if config failed, by default use test mod
	if err != nil {
		config = Config{
			DbUserName: "root",
			DbPassword: "root",
			DbName:     "testdb",
			DbIP:       "127.0.0.1:3306",
		}
	} else {
		ByteData, _ := ioutil.ReadAll(JsonConfig)
		err = json.Unmarshal(ByteData, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}
