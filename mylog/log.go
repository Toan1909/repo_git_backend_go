package mylog
import (
	"runtime"
	log "github.com/Sirupsen/logrus"
	// "github.com/lestrrat/go-file-rotatelogs"
	// "github.com/rifflock/lfshook"
)
func LogError(err error) {
    if err != nil {
        // Retrieve the caller information
        _, file, line, ok := runtime.Caller(1)
        if ok {
            log.Printf("Error: %v (file: %s, line: %d)\n", err, file, line)
        } else {
            log.Printf("Error: %v\n", err)
        }
    }
}