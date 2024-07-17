package nw

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/northwind/internal/enums/appenv"
)

// general contains the general application config
type general struct {
	AppName string
	Env     appenv.Enum
	Debug   bool
}

// api contains the API config
type api struct {
	Port string
}

// ui contains the config for the UI which accesses the API
type ui struct {
	Url string
}

// Config contains all configuration settings
type Config struct {
	General      general
	Db           lyspgdb.Database `toml:"database"`
	DbSuperUser  lyspgdb.User
	DbOwnerUser  lyspgdb.User
	DbServerUser lyspgdb.User
	DbCliUser    lyspgdb.User
	API          api
	UI           ui
}

func (c *Config) LoadFromFile(configFilePath string) (err error) {

	// ensure supplied path exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("configFilePath does not exist: %s", configFilePath)
	} else if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}

	// read conf from toml file
	if _, err := toml.DecodeFile(configFilePath, c); err != nil {
		return fmt.Errorf("toml.DecodeFile failed: %w", err)
	}

	return nil
}
