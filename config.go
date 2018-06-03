package danmu

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/larspensjo/config"
	"os"
	"runtime"
)

var (
	Conf        *Config
	appPath = os.Getenv("GOPATH") + "/src/github.com/kong36088/danmu/"
	configFile  = flag.String("config", appPath+"config/config.ini", "General configuration file")
)

type Config struct {
	values map[string]map[string]string
}

//topic list

func NewConfig() *Config {
	return &Config{
		values: make(map[string]map[string]string),
	}
}

func InitConfig() error {
	Conf = NewConfig()

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	cfgSecs, err := config.ReadDefault(*configFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to find %s %s", *configFile, err))
	}

	for _, section := range cfgSecs.Sections() {
		options, err := cfgSecs.SectionOptions(section)
		if err != nil {
			log.Error("Read options of file %s section %s  failed, %s\n", *configFile, section, err)
			continue
		}
		Conf.values[section] = make(map[string]string)
		for _, v := range options {
			option, err := cfgSecs.String(section, v)
			if err != nil {
				log.Error("Read file %s option %s failed, %s\n", *configFile, v, err)
				continue
			}
			Conf.values[section][v] = option
		}
	}
	return nil
}

func (c *Config) GetConfig(section, option string) string {
	return c.values[section][option]
}

func (c *Config) GetSectionConfig(section string) map[string]string {
	return c.values[section]
}

func (c *Config) GetAllConfig() map[string]map[string]string {
	return c.values
}
