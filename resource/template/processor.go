package template

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/kelseyhightower/confd/log"
)

type Processor interface {
	Processor()
}

func getTemplateResources(config Config)

func Process(config Config) error {
	ts, err := getTe
}
