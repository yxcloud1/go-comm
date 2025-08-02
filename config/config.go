package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

type Configure interface {
	SetDefault() error
}

var (
	workpath   string
	path       string
	configFile string
)

func WorkPath() string {
	return workpath
}

func init() {
	path = *flag.String("path", "", "配置文件路径")
	if path != "" {
		if filepath.IsAbs(workpath){
			workpath = path
		}else{
			workpath, _ = filepath.Abs(os.Args[0])
		}
	}else{
		workpath, _ = filepath.Abs(os.Args[0])
		workpath = filepath.Join(workpath, path)
	}
	configFile = filepath.Join(workpath, "config.json")
	_, err := os.Stat(configFile)
	if err != nil {
		configFile = filepath.Join(workpath,  "config.json")
		if _, err = os.Stat(configFile); err != nil {
			path = filepath.Join(workpath, "config.json")
		}
	}
}

func Load(cfg interface{}) error {
	buf, err := os.ReadFile(path)
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
	if t, ok := cfg.(Configure); ok {
		t.SetDefault()
	}
	json, _ := json.MarshalIndent(cfg, "\t", "\t")
	return os.WriteFile(path, json, 0644)
}

func Save(cfg interface{}) error {
	json, _ := json.MarshalIndent(cfg, "\t", "\t")
	return os.WriteFile(path, json, 0644)
}
