package logger

import (
	"fmt"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	var message[] interface{}
	message = append(message, "A")
	message = append(message, "B")
	log.Println(fmt.Sprint(message...))
}