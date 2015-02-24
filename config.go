package main

import ()

var (
	configFile        = ""
	defaultConfigFile = "/etc/registerd/registerd.toml"
	clientCaKeys      string
	clientCert        string
	clientKey         string
	confdir           string
	config            Config // holds the global registerd config.
	debug             bool
	nodes             Nodes
	printVersion      bool
	quiet             bool
	scheme            string
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
	Interval     int      `toml:"interval"`
	Quiet        bool     `toml:"quiet"`
	Scheme       string   `toml:"scheme"`
	Verbose      bool     `toml:"verbose"`
	Watch        bool     `toml:"watch"`
}
