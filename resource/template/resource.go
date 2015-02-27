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
type ConfigurationConfig struct {
	TemplateResource TemplateResource `toml:"configuration"`
}

// TemplateResource is the representation of a parsed template resource

type configGrouping struct {
	Binary        string
	configOptions map[string]configOption
}

type configOption struct {
	Key   string
	Value string
}
