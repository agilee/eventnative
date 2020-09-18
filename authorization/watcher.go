package authorization

import (
	"crypto/md5"
	"fmt"
	"github.com/ksensehq/eventnative/appstatus"
	"log"
	"time"
)

func Watcher(source string, loadFunc func(string) ([]byte, error), updateFunc func(map[string][]string), reloadSec int) (map[string][]string, error) {
	payload, err := loadFunc(source)
	if err != nil {
		return nil, err
	}

	go watch(source, payload, loadFunc, updateFunc, reloadSec)
	return parseFromBytes(source, payload)
}

func watch(source string, firstTimePayload []byte, loadFunc func(string) ([]byte, error), updateFunc func(map[string][]string), reloadSec int) {
	hash := fmt.Sprintf("%x", md5.Sum(firstTimePayload))
	for {
		if appstatus.Instance.Idle {
			break
		}

		time.Sleep(time.Duration(reloadSec) * time.Second)
		actualPayload, err := loadFunc(source)
		if err != nil {
			log.Printf("Error reloading %s: %v", source, err)
			continue
		}

		newHash := fmt.Sprintf("%x", md5.Sum(actualPayload))
		if hash != newHash {
			result, err := parseFromBytes(source, actualPayload)
			if err != nil {
				log.Printf("Error parsing reloaded %s: %v", source, err)
				continue
			}

			updateFunc(result)
			log.Printf("New resource from %s was loaded", source)
			hash = newHash
		}
	}
}