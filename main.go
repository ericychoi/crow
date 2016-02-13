package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	fsnotify "gopkg.in/fsnotify.v1"
)

const (
	CREATE = iota
	WRITE
	CHMOD
)

var mu *sync.Mutex
var fileStateMap map[string]int

var isVerbose bool = false

// watches for a new file to show up in a given directory, then return the filename
func main() {
	log.SetOutput(os.Stderr)
	if len(os.Args) < 2 {
		log.Fatalf("usage: crow DIR_TO_WATCH")
	}
	dir := os.Args[1]
	mu = &sync.Mutex{}
	fileStateMap = make(map[string]int, 10)
	done := make(chan bool)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if isVerbose {
		log.Printf("start watching %s...\n", dir)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if isVerbose {
					log.Println("event:", event)
				}

				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					if isVerbose {
						log.Printf("created file: %s\n", event.Name)
					}
					mu.Lock()
					fileStateMap[event.Name] = CREATE
					mu.Unlock()

				case event.Op&fsnotify.Write == fsnotify.Write:
					if isVerbose {
						log.Printf("modified file: %s\n", event.Name)
					}

					mu.Lock()
					state, _ := fileStateMap[event.Name]
					if state == CREATE {
						fileStateMap[event.Name] = WRITE
					}
					mu.Unlock()

				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					if isVerbose {
						log.Printf("modified file: %s\n", event.Name)
					}

					mu.Lock()
					state, _ := fileStateMap[event.Name]
					if state == WRITE {
						delete(fileStateMap, event.Name)
					}
					fmt.Fprintln(os.Stdout, event.Name)
					os.Exit(0)
					// don't anybody else come in and kill the script
					//mu.Unlock()
				}

			case err := <-watcher.Errors:
				if err != nil {
					log.Println("error:", err)
				}
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
