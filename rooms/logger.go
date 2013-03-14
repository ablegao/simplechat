package main

import (
	"log"
	"os"
)

//log.Llongfile
var Log = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
