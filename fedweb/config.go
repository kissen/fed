package main

import (
	"flag"
	"log"
	"net/url"
	"sync"
)

type FedWebConfig struct {
	Base *url.URL
}

var cache struct {
	config *FedWebConfig
	lock   sync.Mutex
}

func Config() *FedWebConfig {
	// only initialize once

	cache.lock.Lock()
	defer cache.lock.Unlock()

	// create if necessary

	if cache.config == nil {
		// read in args

		bptr := flag.String("base", "", "base address")

		flag.Parse()

		// evaluate what we got

		if bptr == nil || len(*bptr) == 0 {
			log.Fatal("missing -base argument")
		}

		burl, err := url.Parse(*bptr)
		if err != nil {
			log.Fatal("bad -base argument:", err)
		}

		// return the struct

		cache.config = &FedWebConfig{
			Base: burl,
		}
	}

	// return config

	return cache.config
}
