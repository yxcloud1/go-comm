package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type SqlConnecton struct {
	ConnectionString string `json:"url"`
	ConnectionType   string `json:"type"`
}

type API struct {
	Group string `json:"group"`
}

type LabelStyle struct {
	URL     string `json:"url"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	PageNum int    `json:"pageNum"`
	DPI     int    `json:"dpi"`
}

type UdpDataCollection struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type DataCollection struct {
	Bar         UdpDataCollection `json:"bar"`
	Weight1     UdpDataCollection `json:"weight1"`
	Weight2     UdpDataCollection `json:"weight2"`
	WNo         int               `json:"wdno"`
	MoistureUrl string            `json:"moistureUrl"`
}

type PktApi struct {
	Api struct {
		SendProductionInfo string `json:"url_send_production_info"`
	} `json:"api"`
}

type Config struct {
	API              API            `json:"api"`
	SqlConnecton     SqlConnecton   `json:"db"`
	PktApi           PktApi         `json:"pktapi"`
	DataCollection   DataCollection `json:"datacollection"`
	ProductCertLabel LabelStyle     `json:"productCertLabel"`
	CustomLabel      LabelStyle     `json:"customLabelStyle"`
	ProductLabel     LabelStyle     `json:"productLabel"`
}

var conf *Config

func Conf() *Config {
	if conf == nil {
		loadConfig()
	}
	return conf
}

func loadConfig() {
	log.Println("package Config init()")
	conf = &Config{}
	path := filepath.Dir(os.Args[0]) + "/config.json"
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		path = filepath.Dir(os.Args[0]) + "/../config.json"
		buf, err = ioutil.ReadFile(path)
	}
	if err != nil {
		log.Println(err)
	} else {
		err = json.Unmarshal(buf, conf)
		if err != nil {
			log.Println("json error ", err)
		}
	}
	json, _ := json.MarshalIndent(conf, "\t", "\t")
	ioutil.WriteFile(path, json, 0644)
}
