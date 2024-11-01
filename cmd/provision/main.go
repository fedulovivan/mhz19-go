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
	host = "localhost:7070"
)

func main() {
	fmt.Println("Hello from provisioning tool")
	if value, ok := os.LookupEnv("HOST"); ok {
		host = value
	}
	if dir, ok := os.LookupEnv("DIR"); ok {
		processDir(dir)
	} else {
		fmt.Println("not enough options..\nusage examples:\nDIR=devices make provision\nDIR=rules/system make provision\nDIR=rules/user make provision\nHOST=macmini:7070 DIR=rules/system make provision")
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%stook %s\n", name, elapsed)
}

func processDir(dir string) {
	defer timeTrack(time.Now(), "")
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
		go func(e fs.DirEntry) {
			defer wg.Done()
			fmt.Printf("%s IsDir=%v\n", e.Name(), e.IsDir())
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
				filePath := fmt.Sprintf("%s/%s", dirPath, e.Name())
				res, err := processFile(filePath, apiEntityPath)
				if err == nil {
					succeeded.Add(1)
					fmt.Println("success:", res)
				} else {
					failed.Add(1)
					fmt.Println("fail:", err.Error())
				}
			}
		}(e)
	}
	wg.Wait()
	fmt.Println("---")
	fmt.Println("succeeded", succeeded.Load())
	fmt.Println("failed", failed.Load())
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
