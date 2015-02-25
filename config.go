package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/brendanrjohnson/registerd/backends"
	"github.com/kelseyhightower/confd/log"
	"github.com/kelseyhightower/confd/resource/template"
)

var (
	configFile        = ""
	defaultConfigFile = "/etc/registerd/registerd.toml"
	backend           string
	clientCaKeys      string
	clientCert        string
	clientKey         string
	confdir           string
	config            Config // holds the global registerd config.
	debug             bool
	keepStageFile     bool
	nodes             Nodes
	noop              bool
	onetime           bool
	prefix            string
	printVersion      bool
	quiet             bool
	scheme            string
	templateConfig    template.Config
	backendsConfig    backends.Config
	verbose           bool
	watch             bool
)

// A Config structure is used to configure confd.
type Config struct {
	Backend      string   `toml:"backend"`
	BackendNodes []string `toml:"nodes"`
	ClientCaKeys string   `toml:"client_cakeys"`
	ClientCert   string   `toml:"client_cert"`
	ClientKey    string   `toml:"client_key"`
	ConfDir      string   `toml:"confdir"`
	Debug        bool     `toml:"debug"`
	Noop         bool     `toml:"noop"`
	Prefix       string   `toml:"prefix"`
	Quiet        bool     `toml:"quiet"`
	Scheme       string   `toml:"scheme"`
	Verbose      bool     `toml:"verbose"`
	Watch        bool     `toml:"watch"`
}

func init() {
	flag.StringVar(&backend, "backend", "etcd", "backend to use")
	flag.StringVar(&clientCaKeys, "client-ca-keys", "", "client ca keys")
	flag.StringVar(&clientCert, "client-cert", "", "the client cert")
	flag.StringVar(&clientKey, "client-key", "", "the client key")
	flag.StringVar(&confdir, "confdir", "/etc/confd", "confd conf directory")
	flag.StringVar(&configFile, "config-file", "", "the confd config file")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&keepStageFile, "keep-stage-file", false, "keep staged files")
	flag.Var(&nodes, "node", "list of backend nodes")
	flag.BoolVar(&noop, "noop", false, "only show pending changes")
	flag.BoolVar(&onetime, "onetime", false, "run once and exit")
	flag.StringVar(&prefix, "prefix", "/", "key path prefix")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.BoolVar(&quiet, "quiet", false, "enable quiet logging")
	flag.StringVar(&scheme, "scheme", "http", "the backend URI scheme (http or https)")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&watch, "watch", false, "enable watch support")
}

// initConfig initializes the confd configuration by first setting defaults,
// then overriding setting from the confd config file, and finally overriding
// settings from flags set on the command line.
// It returns an error if any.
func initConfig() error {
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
	}
	// Set defaults.
	config = Config{
		ConfDir: "/etc/registerd",
		Scheme:  "http",
	}

	// Update config from the TOML configuration file
	if configFile == "" {
		log.Warning("Skipping registerd config file.")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}
	// Update config from commandline flags
	processFlags()

	// Configure logging
	log.SetQuiet(config.Quiet)
	log.SetVerbose(config.Verbose)
	log.SetDebug(config.Debug)

	// Update BackendNodes
	if len(config.BackendNodes) == 0 {
		switch config.Backend {
		case "consul":
			config.BackendNodes = []string{"127.0.0.1:8500"}
		case "etcd":
			peerstr := os.Getenv("ETCDCTL_PEERS")
			if len(peerstr) > 0 {
				config.BackendNodes = strings.Split(peerstr, ",")
			} else {
				config.BackendNodes = []string{"http://127.0.0.1:4001"}
			}
		}
	}

	// Initialize the storage client
	log.Notice("Backend set to " + config.Backend)

	//Backend configuration.
	backendsConfig = backends.Config{
		Backend:      config.Backend,
		ClientCaKeys: config.ClientCaKeys,
		ClientCert:   config.ClientCert,
		ClientKey:    config.ClientKey,
		BackendNodes: config.BackendNodes,
		Scheme:       config.Scheme,
	}

	// Service Definition Configuration
	templateConfig = template.Config{
		ConfDir:       config.ConfDir,
		ConfigDir:     filepath.Join(config.ConfDir, "conf.d"),
		KeepStageFile: keepStageFile,
		Noop:          config.Noop,
		Prefix:        config.Prefix,
	}
	return nil
}

func processFlags() {
	flag.Visit(setConfigFromFlag)
}

func setConfigFromFlag(f *flag.Flag) {
	switch f.Name {
	case "debug":
		config.Debug = debug
	case "client-cert":
		config.ClientCert = clientCert
	case "client-key":
		config.ClientKey = clientKey
	case "client-cakeys":
		config.ClientCaKeys = clientCaKeys
	case "confdir":
		config.ConfDir = confdir
	case "node":
		config.BackendNodes = nodes
	case "quiet":
		config.Quiet = quiet
	case "scheme":
		config.Scheme = scheme
	case "verbose":
		config.Verbose = verbose
	case "watch":
		config.Watch = watch
	}
}
