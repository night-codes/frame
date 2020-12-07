package frame

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   = log.New(ioutil.Discard, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info    = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error   = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
)
