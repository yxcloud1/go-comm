package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	path    string
	workDir string = ""
)

type Configure interface {
	SetDefault() error
}

func init() {
	path = ""
	for _, v := range os.Args {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "-workdir=") {
			path = v[len("-workdir="):]
			break
		}
	}
	if filepath.IsAbs(path) {
		workDir = path
	} else {
		workDir, _ = filepath.Abs(os.Args[0])
		workDir = filepath.Dir(workDir)
		workDir = filepath.Join(workDir, path)
	}
	if info , err := os.Stat(workDir); err != nil{
		workDir, _ = filepath.Abs(os.Args[0])
		workDir = filepath.Dir(workDir)
	}else{
		if !info.IsDir(){
			workDir = filepath.Dir(workDir)
		}
	}
}

func WorkDir() string {
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
	return configFile
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
