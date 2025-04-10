package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/yummynaan/kit-enrollment-helper/internal/app/task"
	"github.com/yummynaan/kit-enrollment-helper/internal/service"
)

var (
	isArchive    bool
	year         int
	page         int
	targetDomain string
	baseURL      string
)

func init() {
	var fy int
	start := time.Now()
	m := start.Month()
	if m >= time.April {
		fy = start.Year()
	} else {
		fy = start.Year() - 1
	}

	flag.BoolVar(&isArchive, "archive", false, "whether to retrieve info from the archive")
	flag.IntVar(&year, "year", fy, "fiscal year")
	flag.IntVar(&page, "page", 1, "page offset")
	flag.Parse()

	targetDomain = "www.syllabus.kit.ac.jp"
	baseURL = "https://" + targetDomain + "/"
	if isArchive {
		baseURL += "archive/"
	}
}

func main() {
	param := "?c=search_list&sk=&dc=01" + fmt.Sprintf("&yr=%d&page=%d", year, page)
	targetURL, err := url.Parse(baseURL + param)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := service.CreateRepository()
	if err != nil {
		log.Fatal(err)
	}

	worker := task.NewFetchSyllabusWorker(targetURL, year, repository)
	if err := worker.Run(); err != nil {
		log.Fatal(err)
	}
}
