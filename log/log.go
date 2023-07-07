package log

import (
	"fmt"
	"log"
)

func LogMe(name, message string) {
	log.Println(fmt.Sprintf("[%s] - %s", name, message))
}
