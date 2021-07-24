package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	DatePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	DateLayout  = "2006-01-02"
)

func main() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("exited with error:", (*err).Error())
			os.Exit(1)
		} else {
			log.Println("exited")
		}
	}(&err)

	flag.Parse()

redo:
	now := time.Now()
	for _, dir := range flag.Args() {
		if err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return nil
			}
			handleFile(path, info, now)
			return nil
		}); err != nil {
			return
		}
	}

	log.Println("SLEEPING")
	time.Sleep(time.Hour)

	goto redo
}

func handleFile(path string, info fs.FileInfo, now time.Time) {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("ERROR:", (*err).Error())
		}
	}(&err)

	log.Println("FOUND:", path)

	match := DatePattern.FindString(info.Name())
	if match == "" {
		if info.Size() > 5*1000*1000*1000 {
			log.Println("> TRUNCATE")
			if err = os.Truncate(path, 0); err != nil {
				return
			}
		} else {
			log.Println("> KEEP")
		}
	} else {
		var date time.Time
		if date, err = time.Parse(DateLayout, match); err != nil {
			return
		}
		if now.Sub(date) > time.Hour*24*4 {
			log.Println("> DELETE")
			if err = os.Remove(path); err != nil {
				return
			}
		} else {
			log.Println("> KEEP")
		}
	}

}
