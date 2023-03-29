package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)
var
(
	path string
)

func init(){
	path = filepath.Dir(os.Args[0]) + "/config.json"
	_, err := os.Stat(path)
	if err != nil{
		path = filepath.Dir(os.Args[0]) + "/../config.json"
		if _, err = os.Stat(path);err != nil{
			path  = filepath.Dir(os.Args[0]) + "/config.json"
		}
	}
}

func Load(cfg interface{})error{
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	} else {
		err = json.Unmarshal(buf, cfg)
		if err != nil {
			log.Println("json error ", err)
			return err
		}
	}
	json, _ := json.MarshalIndent(cfg, "\t", "\t")
	return ioutil.WriteFile(path, json, 0644)
}

func Save(cfg interface{}) error{
	json, _ := json.MarshalIndent(cfg, "\t", "\t")
	return ioutil.WriteFile(path, json, 0644)
}
