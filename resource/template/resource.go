package template

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/brendanrjohnson/registerd/backends"
	"github.com/kelseyhightower/confd/log"
	"github.com/kelseyhightower/memkv"
)

type Config struct {
	ConfDir       string
	ConfigDir     string
	KeepStageFile bool
	Noop          bool
	Prefix        string
	StoreClient   backends.StoreClient
	Template      string
}

// ConfigurationConfig holds the parsed configuration resource.
type ConfigurationResoureConfig struct {
	ConfigurationResource ConfigurationResource `toml:"configuration"`
}

// ConfigurationResource is the representation of a parsed template resource

type ConfigurationResource struct {
	Binary			string
	ConfigOptions 	map[string]configOption
	funcMap			map[string]interface{}
	keepStageFile	bool
	noop 			bool
	prefix 			string
	store			memkv.Store
	storeClient		backends.StoreClient
	}
}

type configOption struct {
	Key   string
	Value string
}

var ErrEmptySrc = errors.New("empty src template")

// NewConfigurationResource creates a ConfigurationResource
func NewConfigurationResource(path string, config Config) (*ConfigurationResource, error) {
	if config.StoreClient == nil {
		return nil, errors.New("A valid StoreClient is required.")
	}
	var cc *ConfigurationResoureConfig
	log.Debug("Loading template resource from" + path)
	_, err := toml.DecodeFile(path, &cc)
	if err != nil {
		return nil, fmt.Errorf("Cannot process template resource %s - %s", path, err.Error())
	}
	cr := cc.ConfigurationResource
	cr.keepStageFile = config.KeepStageFile
	cr.noop = config.Noop
	cr.storeClient = config.StoreClient
	cr.funcMap = newFuncMap()
	cr.store = memkv.New()
	addFuncs(cr.funcMap, cr.store.FuncMap)
	cr.prefix = filepath.Join("/", config.Prefix, cr.prefix)
	return &cr, nil
}

// process is a convenience function that wraps calls to the three main tasks
// required to keep local configuration files in sync. First we gather vars
// from the store, then we stage a candidate configuration file, and finally sync
// things up.
// It returns an error if any.
func (c *ConfigurationResource) process() error {
	if err := c.Set
}