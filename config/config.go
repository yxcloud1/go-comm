package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)
var(
	path string
)
type Configure interface {
	SetDefault() error
}

func init(){
	path = ""
	flag.StringVar(&path, "workdir", "", "配置文件路径")
}

func WorkDir() string {
	workDir := ""
	if path != "" {
		info, err := os.Stat(path)
		if err != nil{
			path = ""
		}else if info.IsDir(){
			path = ""
		}
	}
	if path != ""{
		if filepath.IsAbs(path) {
			workDir = path
		} else {
			workDir, _ = filepath.Abs(os.Args[0])
			workDir = filepath.Join(workDir, path)
		}
	} else {
		workDir, _ = filepath.Abs(os.Args[0])
	}
	return workDir
}

func configFile() string {
	workDir := WorkDir()
	configFile := filepath.Join(workDir, "config.json")
	_, err := os.Stat(configFile)
	if err != nil {
		configFile = filepath.Join(workDir, "config.json")
		if _, err = os.Stat(configFile); err != nil {
			configFile = filepath.Join(workDir, "config.json")
		}
	}
	return  configFile
}

func Load(cfg interface{}) error {
	buf, err := os.ReadFile(configFile())
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
	return os.WriteFile(configFile(), json, 0644)
}

func Save(cfg interface{}) error {
	json, _ := json.MarshalIndent(cfg, "\t", "\t")
	return os.WriteFile(configFile(), json, 0644)
}
