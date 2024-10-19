package config

import (
	"log"
	"time"
)

type Stopper interface {
	Stop()
}

func startWatch(cfgFilePath string, modTime time.Time) Stopper {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			modTime = watchFile(cfgFilePath, modTime)
		}
	}()

	return ticker
}

func watchFile(cfgFilePath string, modTime time.Time) time.Time {
	newModTime, err := getModTime(cfgFilePath)
	if err != nil {
		log.Printf("failed to stat config file: %+v\n", err)
		return modTime
	}

	if newModTime.After(modTime) {
		newCfg, err := readConfigFile(cfgFilePath)
		if err != nil {
			log.Printf("failed to read config file: %+v\n", err)
			return modTime
		}
		setV(newCfg)
		log.Println("config file updated")
		return newModTime
	}

	return modTime
}
