package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const cacheRootFolder = "./data/cache/"
const cacheExtension = ".gob"

func getCacheFilename(cacheName string) string {
	return cacheRootFolder + cacheName + cacheExtension
}

func LoadFromCache(cacheName string, outputCachedData bool, cacheItem interface{}) (cacheMiss bool) {
	cacheFileAsBytes, err := ioutil.ReadFile(getCacheFilename(cacheName))
	cacheMiss = true
	if err == nil {
		// Decode from cache, ignore failures
		err = gob.NewDecoder(bytes.NewBuffer(cacheFileAsBytes)).Decode(cacheItem)
		if err == nil {
			cacheMiss = false
			if outputCachedData {
				fmt.Printf("Loaded cache '%v' (%v bytes) with: %+v\n", cacheName, len(cacheFileAsBytes), cacheItem)
			} else {
				fmt.Printf("Loaded cache '%v' (%v bytes).\n", cacheName, len(cacheFileAsBytes))
			}
		} else {
			fmt.Printf("Could not decode cache '%v': %v.\n", cacheName, err)
		}
	} else {
		fmt.Printf("Could not read cache '%v': %v.\n", cacheName, err)
	}

	return cacheMiss
}

func SaveToCache(cacheName string, cacheItem interface{}) bool {
	byteBuffer := new(bytes.Buffer)
	err := gob.NewEncoder(byteBuffer).Encode(cacheItem)
	if err != nil {
		fmt.Printf("Unable to encode the '%v' cache: %v\n", cacheName, err)
		return false
	} else {
		err = ioutil.WriteFile(getCacheFilename(cacheName), byteBuffer.Bytes(), os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to cache the '%v' encoded data: %v\n", cacheName, err)
			return false
		}
	}

	return true
}
