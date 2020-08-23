package main

import (
	"flag"
	"net/http"
	"orealtest/config"
	"orealtest/service"
)

func init() {
	flag.StringVar(&config.Root, "P", "/", "Root path")
}

func main() {
	flag.Parse()
	http.HandleFunc("/dir", service.DirHandler)
	http.HandleFunc("/dirInfo", service.DirInfoHandlder)
	http.ListenAndServe(":9000", nil)
}
