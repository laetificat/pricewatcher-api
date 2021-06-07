package log

import "log"

func main() {

}

func Info(m interface{}) {
	log.Printf("INFO %v", m)
}

func Debug(m interface{}) {
	log.Printf("DEBUG %v", m)
}

func Warning(m interface{}) {
	log.Printf("WARNING %v", m)
}

func Error(m interface{}) {
	log.Printf("ERROR %v", m)
}

func Panic(m interface{}) {
	log.Printf("PANIC %v", m)
}

func Fatal(m interface{}) {
	log.Printf("FATAL %v", m)
}
