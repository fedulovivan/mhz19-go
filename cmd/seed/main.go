package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	basePath = "./assets"
)

var (
	host   = "localhost:7070"
	filter = ""
)

var help = `
not enough options...
usage examples:
  DIR=devices make seed
  DIR=rules/system make seed
  DIR=rules/user make seed
  HOST=%s DIR=rules/system make seed
  DIR=rules/system FILTER=buried make seed
`

func main() {
	fmt.Println("Hello from data seeding tool")
	if v, ok := os.LookupEnv("HOST"); ok {
		host = v
	}
	if v, ok := os.LookupEnv("FILTER"); ok {
		filter = v
	}
	if dir, ok := os.LookupEnv("DIR"); ok {
		processDir(dir)
	} else {
		fmt.Printf(help, host)
	}
}

// use simplified version of github.com/fedulovivan/mhz19-go/pkg/utils::TimeTrack()
// to avoid extra work with slog initialisation
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%stook %s\n", name, elapsed)
}

func processDir(dir string) {
	defer TimeTrack(time.Now(), "")
	succeeded := atomic.Int32{}
	failed := atomic.Int32{}
	dirPath := path.Join(basePath, dir)
	entries, err := os.ReadDir(dirPath)
	dd := strings.Split(dir, "/")
	apiEntityPath := dd[0]
	if err != nil {
		panic(err.Error())
	}
	var wg sync.WaitGroup
	for _, e := range entries {
		wg.Add(1)
		go func(entry fs.DirEntry) {
			defer wg.Done()
			name := entry.Name()
			fmt.Printf("%s IsDir=%v\n", name, entry.IsDir())
			use := !entry.IsDir() &&
				strings.HasSuffix(name, ".json") &&
				(len(filter) == 0 || strings.Contains(name, filter))
			if use {
				filePath := fmt.Sprintf("%s/%s", dirPath, name)
				res, err := processFile(filePath, apiEntityPath)
				if err == nil {
					succeeded.Add(1)
					fmt.Println("✅ success:", name, res)
				} else {
					failed.Add(1)
					fmt.Println("❌ fail:", name, err.Error())
				}
			} else {
				fmt.Println("⏭️  skipped:", name)
			}
		}(e)
	}
	wg.Wait()
	fmt.Println("---")
	fmt.Println("succeeded", succeeded.Load())
	fmt.Println("failed", failed.Load())
	fmt.Println("---")
	fmt.Println("host", host)
	fmt.Println("dir", dir)
	fmt.Println("---")
}

func processFile(filePath string, apiEntityPath string) (res string, err error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	url := fmt.Sprintf(
		"http://%v/api/%s",
		host,
		apiEntityPath,
	)
	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(file),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	var bodyRaw []byte
	var bodyParsed map[string]any
	bodyRaw, err = io.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyRaw, &bodyParsed)
	if err != nil {
		return
	}
	if v, ok := bodyParsed["is_error"]; ok {
		if vv, ok := v.(bool); ok && vv {
			err = fmt.Errorf(bodyParsed["error"].(string))
			return
		}
	}
	res = strings.Trim(string(bodyRaw), "\n")
	return
}
