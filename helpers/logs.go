package helpers

import (
	"fmt"
	"log"
)

func LogError(v ...any) {
	errorMessage := fmt.Sprintf("\x1b[41m\x1b[37m  ERROR  \x1b[0m %v", v)
	log.Println(errorMessage)
}

func LogFatal(v ...any) {
	fatalMessage := fmt.Sprintf("\x1b[41m\x1b[37m  FATAL!   APP WLL EXIT NOW \x1b[0m %v", v)
	log.Fatal(fatalMessage)
}

func LogInfo(v ...any) {
	infoMessage := fmt.Sprintf("\x1b[47m\x1b[30m  INFO   \x1b[0m %v", v)
	log.Println(infoMessage)
}

func LogOK(v ...any) {
	okMessage := fmt.Sprintf("\x1b[42m\x1b[37m   OK    \x1b[0m %v", v)
	log.Println(okMessage)
}

func LogWarning(v ...any) {
	warningMessage := fmt.Sprintf("\x1b[43m\x1b[30m WARNING \x1b[0m %v", v)
	log.Println(warningMessage)
}
