package main

import (
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func watchFiles(dir, userID string, cfg *Config) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Create a channel to receive OS signals CTRL+C
	sigs := make(chan os.Signal, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at current working directory. Dispatch remove events to c.
	if err := notify.Watch(dir, c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	go func() {
		for {
			select {
			case ei := <-c:
				//_ = ei
				//log.Println("Got event:", ei.Event(), "on file:", ei.Path())
				if ei.Event() == notify.Create {
					// check if path is a directory
					// since it was just created we don't upload empty directories
					fi, err := os.Stat(ei.Path())
					if err != nil {
						log.Println(err)
						continue
					}
					if fi.IsDir() {
						continue
					}
				}

				files := getFiles(dir)
				uploadFiles(userID, files, cfg)
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a CTRL+C signal is received
	<-sigs
}
