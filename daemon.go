package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
)

func setupNotify(filename string) (chan notify.EventInfo, string) {
	absfn, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Cannot convert %v to absolute path: %v", filename, err)
	}

	dir := filepath.Dir(absfn)
	basename := filepath.Base(absfn)

	cn := make(chan notify.EventInfo, 1)
	err = notify.Watch(dir, cn, notify.All)
	if err != nil {
		log.Fatal("Cannot setup watch for ", dir, ": ", err)
	}

	log.Infof("Watching %v for changes to %v", dir, basename)

	return cn, absfn
}

func runIfInputExists(filename string) {
	if _, err := os.Stat(filename); err != nil {
		log.Error("Cannot stat hosts file, skipping sync: ", err)
	} else {
		runOnce()
	}
}

func daemon(interval time.Duration, filename string) {
	cn, absfn := setupNotify(filename)
	defer notify.Stop(cn)

	log.Info("Running initial sync")
	runIfInputExists(absfn)

	log.Info("sync scheduled every ", interval)
	ticker := time.NewTicker(interval)

	for {
		resyncNeeded := false
		// Block on either the ticker or inotify
		select {
		case <-ticker.C:
			resyncNeeded = true
		case ei := <-cn:
			if ei.Path() == absfn {
				log.Info("file change event detected: ", ei)
				resyncNeeded = true
			}
		}

		if resyncNeeded {
			runIfInputExists(absfn)
		}
	}
}
