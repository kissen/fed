package config

import (
	"github.com/BurntSushi/toml"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
)

type FedConfig struct {
	// Hostname under which the instance is reachable in the
	// open web. Something like "fed.example.com"
	Hostname string

	// What address to listen on. Something like "localhost:8080"
	// if you are using a reverse proxy like nginx.
	ListenAddress string

	// Location of the storage file. The process running fed
	// will need rw permissions on that file and the directory
	// it is in.
	StorageFile string
}

// Pointer to the singelton instance of the global config.
var singleton *FedConfig
var once sync.Once

// Return the singleton instance of the configuration.
func Get() *FedConfig {
	once.Do(fillInSingleton)
	return singleton
}

// Return the Hostname property as a newly created URL. This is the
// base URL under which the instance should be reachable on the open
// internet.
func (fc *FedConfig) GlobalURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   fc.Hostname,
	}
}

// Fill in the singleton global or stop the program on failure.
func fillInSingleton() {
	filename := "doc/fed.conf"

	if c, err := loadConfigFrom(filename); err != nil {
		log.Println(err)
		log.Fatal("cannot start without configuration file")
	} else {
		singleton = c
	}
}

// Try to open and parse configuration file at filename.
func loadConfigFrom(filename string) (*FedConfig, error) {
	var c FedConfig

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not load config file")
	}

	if err := toml.Unmarshal(bs, &c); err != nil {
		return nil, errors.Wrapf(err, `filename="%v" not a valid config file`, filename)
	}

	return &c, nil
}
