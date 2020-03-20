package main

import (
	"flag"
	"log"
	"net/url"
	"sync"
)

type FedConfig struct {
	Base  *url.URL
	Store string
}

var cache struct {
	config *FedConfig
	lock   sync.Mutex
}

func Config() *FedConfig {
	// only initialize once

	cache.lock.Lock()
	defer cache.lock.Unlock()

	// create if necessary

	if cache.config == nil {
		// read in args

		bptr := flag.String("base", "", "base address")
		sptr := flag.String("store", "", "storage file location")

		flag.Parse()

		// evaluate what we got

		if bptr == nil || len(*bptr) == 0 {
			log.Fatal("missing -base argument")
		}

		if sptr == nil || len(*sptr) == 0 {
			log.Fatal("missing -store argument")
		}

		burl, err := url.Parse(*bptr)
		if err != nil {
			log.Fatal("bad -base argument:", err)
		}

		// return the struct

		cache.config = &FedConfig{
			Base:  burl,
			Store: *sptr,
		}
	}

	// return config

	return cache.config
}
