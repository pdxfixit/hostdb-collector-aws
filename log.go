package main

import "log"

func init() {
	// log options
	log.SetFlags(log.Ldate + log.Ltime + log.Lshortfile + log.LUTC)
}
