package log

import (
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

func Info(message string) {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	fmt.Print(Cyan("[" + formattedTime + "] "))
	fmt.Print(Blue("[INFO]: "))
	fmt.Println(message)
}

func Error(message string) {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	fmt.Print(Cyan("[" + formattedTime + "] "))
	fmt.Print(Red("[ERROR]: "))
	fmt.Println(message)
}

func Success(message string) {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	fmt.Print(Cyan("[" + formattedTime + "] "))
	fmt.Print(Green("[SUCCESS]: "))
	fmt.Println(message)
}
