package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func LogError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func Download(rawUrl string) []byte {
	// Determine local path based on URL
	fileUrl, err := url.Parse(rawUrl)
	LogError(err)

	filePath := path.Join("./", fileUrl.Host, fileUrl.Path)

	// Create empty file
	err = os.MkdirAll(filepath.Dir(filePath), 0770)
	LogError(err)

	file, err := os.Create(filePath)
	LogError(err)

	// Get request
	resp, err := http.Get(rawUrl)
	LogError(err)

	// Read body
	body, err := ioutil.ReadAll(resp.Body)
	LogError(err)

	// Save body
	file.Write(body)

	defer file.Close()
	defer resp.Body.Close()

	// Print message
	log.Println("Downloaded " + rawUrl)

	return body
}

var processed map[string]bool = make(map[string]bool)

func RecursiveDownload(url string) {
	if _, ok := processed[url]; ok {
		return
	}
	processed[url] = true

	body := Download(url)

	if strings.HasSuffix(url, "json") {
		for _, v := range regexp.MustCompile(`"(https:\/\/[^"]*)"`).FindAllSubmatch(body, -1) {
			RecursiveDownload(string(v[1]))
		}

		for _, v := range regexp.MustCompile(`"hash": "(\w*)"`).FindAllSubmatch(body, -1) {
			Download("https://resources.download.minecraft.net/" + string(v[1][0:2]) + "/" + string(v[1]))
		}
	}

}

func main() {
	RecursiveDownload("https://piston-meta.mojang.com/mc/game/version_manifest.json")
}
