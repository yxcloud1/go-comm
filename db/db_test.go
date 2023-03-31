package db

import (
	"log"
	"testing"
)

func TestOracle(t *testing.T) {
	dsn := "oracle://mdprau:rau_r0ima_mdp@ds1.cmaa.cf:21521/RAUPR1?TRACE FILE=trace.log"
	SetOption("oracle", dsn)
	if result, err := DB().ExecuteQuery("SELECT * FROM  MDPRAU.MDP_QA_MEREELPROF where rownum < 5");err!=nil{
		log.Println(err)
	}else{
		log.Println(result)
	}
}