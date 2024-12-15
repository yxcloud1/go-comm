package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type option struct {
	level  string
	status map[string]bool
	logDuration time.Duration
}

var (
	logMutex sync.Mutex
	errMutex sync.Mutex
	opt      = &option{
		level: "INFO",
		status: map[string]bool{
			"ERROR":   true,
			"WARRING": true,
			"INFO":    true,
			"DEBUG":   true,
		},
		logDuration: time.Hour * 24 * 30 * 2,


	}
	defaultColor = "\033[0m"
	COLOR        = map[string]string{

		"ERROR":   "\033[31m",
		"WARNING": "\033[33",
		"INFO":    "\033[0m",
		"RESET":   "\033[0m",
		"DEBUG":   "\033[32",
	}
)

func SetOption(level string, opts ... string ) {
	opt.level = level
	for k := range opt.status {
		if strings.Contains(strings.ToUpper(level), k) {
			opt.status[k] = true
		}
	}
	for k, v := range opts {
		switch k{
		case 0:
			if d, err := time.ParseDuration(v); err == nil{
				opt.logDuration = d
			}
		}
	}
}

func deleteOld(logDir string) {
	now := time.Now()
	twoMonthsAgo := now.Add(-1*opt.logDuration)
	twoMonthsAgoStr := twoMonthsAgo.Format("200601") // 格式化为 yyyyMM

	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if info.IsDir() && path != logDir {
			dirName := filepath.Base(path)

			if len(dirName) == 6 {
				if _, err := strconv.Atoi(dirName); err == nil {
					if dirName < twoMonthsAgoStr {
						err := os.RemoveAll(path)
						if err != nil {
							fmt.Printf("Failed to delete directory %s: %v\n", path, err)
						} else {
							fmt.Printf("Deleted old directory: %s\n", path)
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", logDir, err)
	}
}

func createPath(logtype string) (string, error) {
	t := time.Now()
	cwd := filepath.Dir(os.Args[0])
	deleteOld(cwd + "/" + logtype)
	path := t.Format("200601")
	fileName := "/" + t.Format("02") + ".txt"
	path = cwd + "/" + logtype + "/" + path
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path + fileName, nil
}

func createLogParh() (string, error) {
	return createPath("log")
}

func createErrParh() (string, error) {
	return createPath("err")
}

func txtLogger(msg ...interface{}) error {
	message := ""
	if len(msg) == 0{
		return nil
	}
	for k, v := range msg {
		if k > 0{
		message += fmt.Sprintf("\r\n%v", v)
	}
}
	beginstring := fmt.Sprintf("%s    ----------------------------%s-----------------------------------------",
								time.Now().Format(time.RFC3339), msg[0])
	endstring := strings.Repeat("*", (len(beginstring)))
	message = beginstring + message + "\r\n"+  endstring
	fmt.Println(message)
	logMutex.Lock()
	defer logMutex.Unlock()
	if path, err := createLogParh(); err == nil {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		//message = fmt.Sprintf("%s\r\n%s\r\n", beginstring , message)
		write.WriteString(message)
		write.Flush()
		return nil
	} else {
		return err
	}
}


func txtError(msg ...interface{}) error {
	message := ""
	if len(msg) == 0{
		return nil
	}
	for k, v := range msg {
		if k > 0{
		message += fmt.Sprintf("\r\n%v", v)
	}
}
	beginstring := fmt.Sprintf("%s    ----------------------------%s-----------------------------------------",
								time.Now().Format(time.RFC3339), msg[0])
	endstring := strings.Repeat("*", (len(beginstring)))
	message = beginstring + message + "\r\n" + endstring
	fmt.Println(message)
	errMutex.Lock()
	defer errMutex.Unlock()
	if path, err := createErrParh(); err == nil {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		//message = fmt.Sprintf("%s\r\n%s\r\n", beginstring , message)
		write.WriteString(message)
		write.Flush()
		return nil
	} else {
		return err
	}
}


func Log(level string, messages ...interface{}) error {
	message := ""
	for _, v := range messages {
		message += fmt.Sprintf("%+v\r\n", v)
	}
	log.Println(message)

	msg := fmt.Sprintf("%s\t%s\r\n%s\r\n", time.Now().Format(time.RFC3339), level, message)

	if cl, ok := COLOR[level]; ok {
		fmt.Printf("%s%s%s\n", cl, msg, defaultColor)
	} else {
		fmt.Printf("%s%s%s\n", defaultColor, msg, defaultColor)
	}

	if t, ok := opt.status[strings.ToUpper(level)]; ok {
		if !t {
			return nil
		}
	}

	logMutex.Lock()
	defer logMutex.Unlock()
	if path, err := createPath(level); err == nil {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		write.WriteString(msg)
		write.Flush()
		return nil
	} else {
		return err
	}
}

func TxtDebug(message ...interface{}) error {
	if strings.Contains(opt.level, "DEBUG") {
		message = append([]interface{}{"DEBUG"}, message...)
		return txtLogger(message...)
	}
	return nil
}

func TxtLog(message ...interface{}) error {
	if strings.Contains(opt.level, "INFO") {
		message = append([]interface{}{"INFO"}, message...)
		return txtLogger(message...)
	}
	return nil
}

func TxtErr(message ...interface{}) error {
	message = append([]interface{}{"ERROR"}, message...)

	txtLogger(message)
	if !strings.Contains(opt.level, "ERROR") {
		return nil
	}
	return txtError(message...)
}
