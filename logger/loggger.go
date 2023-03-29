package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type option struct{
	level string
}

var (
	logMutex sync.Mutex
	errMutex sync.Mutex
	opt = &option{
		level: "INFO",
	}
)



func SetOption(level string){
	opt.level = level
}

func createLogParh() (string, error) {
	t := time.Now()
	cwd := filepath.Dir(os.Args[0])
	path := t.Format("200601")
	fileName := "/" + t.Format("02") + ".txt"
	path = cwd + "/log/" + path
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path + fileName, nil
}

func createErrParh() (string, error) {
	t := time.Now()
	cwd := filepath.Dir(os.Args[0])
	path := t.Format("200601")
	fileName := "/" + t.Format("02") + ".txt"
	path = cwd + "/err/" + path
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path + fileName, nil
}

func txtLogger(message ...interface{}) error {
	log.Println(message...)
	logMutex.Lock()
	defer logMutex.Unlock()
	if path, err := createLogParh(); err == nil {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		msg := fmt.Sprintf("%s\t%s\r\n", time.Now().Format(time.RFC3339), fmt.Sprint(message...))
		write.WriteString(msg)
		write.Flush()
		return nil
	} else {
		return err
	}
}


func TxtDebug(message... interface{}) error {
	if strings.Contains(opt.level, "DEBUG"){
		message = append([]interface{}{"DEBUG"}, message...)
		return txtLogger(message)
	}
	return nil
}

func TxtLog(message... interface{}) error {
	if strings.Contains(opt.level, "INFO"){
		message = append([]interface{}{"INFO"}, message...)
		return txtLogger( message)
	}
	return nil
}


func TxtErr(message ...interface{}) error {
	message = append([]interface{}{"ERROR"}, message...)
	txtLogger( message)
	if !strings.Contains(opt.level, "ERROR"){
		return nil
	}
	errMutex.Lock()
	defer errMutex.Unlock()
	if path, err := createErrParh(); err == nil {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		msg := fmt.Sprintf("%s\t%s\r\n", time.Now().Format(time.RFC3339), fmt.Sprint(message...))
		write.WriteString(msg)
		write.Flush()
		return nil
	} else {
		return err
	}
}
