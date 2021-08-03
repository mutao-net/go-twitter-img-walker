package main

import (
	. "twitter-img-walker/conf"
	"fmt"
	"net/url"
	"regexp"
	"io"
	"net/http"
	"os"
	"flag"
	"bufio"
	"log"
)
const (
	version="1.0.0"
	filePath="./target.txt"
)
func main()  {
	flag.Parse()
	args := flag.Args()
	path := checkOptions(args)
	run(path)
}

func run (path string) {
	users := getTargets()
	fmt.Println(users)
	for _, user := range users {
		fmt.Println(user)
		check_dir(path, user)
		mediaList := getTweetImg(user)
		fmt.Println(mediaList)
		getImg(path, user, mediaList)
	}
}

func checkOptions(args []string) string{
	if len(args) > 1 {
		fmt.Println("too many arguments usage: main.go {image download path} default ./img")
		os.Exit(1)
	}
	if len(args) == 0 {
		return "./img";
	}
	if len(args) == 1 {
		reg := regexp.MustCompile(`([^\.?\/]+.?)\/$`)
		if (reg.MatchString(args[0])) {
			return args[0]
		} else {
			fmt.Println("specify an absolute path")
		}
		os.Exit(1)
	}
	return args[0]
}
// get Screen name
func getTargets() []string{
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
		os.Exit(1)
	}

	defer f.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}

	return lines
}
// Tweet get image urls
func getTweetImg(screenName string) []string {
	api := InitTwitterApi()
	v := url.Values{}
	fmt.Println(screenName)
	v.Set("screen_name", screenName)
	v.Set("count", "1000")

	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		log.Println(err)
	}
	mediaList := []string{}
	for _, tweet := range tweets {
		for _, media := range tweet.Entities.Media {
			mediaList = append(mediaList, media.Media_url_https)
		}
	}
	return mediaList
}

// check if directory exists
func check_dir (path string, screenName string) {
	if _, err := os.Stat(path + "/" + screenName); os.IsNotExist(err) {
		os.MkdirAll(path + "/" + screenName, 0755)
	}
}

// check if file exsits
func fileExists(path string, filename string) bool {
	_, err := os.Stat(path + "/" + filename)
	return err == nil
}

// image download
func getImg(path string, screenName string, mediaList []string) {
	for _, url := range mediaList {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		reg := regexp.MustCompile(`([^\/]+?)(\.jpg|\.jpeg|\.gif|\.png)$`)
		if !fileExists(path, reg.FindString(url)) {
			file, err := os.Create(path + "/" + screenName + "/" + reg.FindString(url))
			if err != nil {
				panic(err)
			}
			defer file.Close()
		
			io.Copy(file, response.Body)
		}
	}
}