package config

import (
	"flag"
	"log"
	"testing"
)

func TestWorkPath(t* testing.T){
	log.Println("--------")
	flag.Parse()
	t.Log(WorkDir())
}