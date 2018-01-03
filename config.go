package danmu

import (
	"flag"
	"github.com/larspensjo/config"
	"log"
	"runtime"
)

var (
	configFile = flag.String("config", "config/config.ini", "General configuration file")
)

//topic list
var CFG map[string]map[string]string

func ReadConfig() map[string]map[string]string {
	if CFG != nil {
		return CFG
	}

	CFG = make(map[string]map[string]string)

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find %s %s", *configFile, err)
	}

	for _, section := range cfg.Sections() {
		options, err := cfg.SectionOptions(section)
		if err != nil {
			log.Printf("Read options of file %s section %s  failed, %s\n", *configFile, section, err)
			continue
		}
		CFG[section] = make(map[string]string)
		for _, v := range options {
			option, err := cfg.String(section, v)
			if err != nil {
				log.Printf("Read file %s option %s failed, %s\n", *configFile, v, err)
				continue
			}
			CFG[section][v] = option
		}
	}

	return CFG
}

func GetConfig(section, option string) string {
	if CFG == nil {
		ReadConfig()
	}
	return CFG[section][option]
}

func GetSectionConfig(section string) map[string]string {
	return CFG[section]
}

func GetAllConfig() map[string]map[string]string {
	if CFG == nil {
		ReadConfig()
	}
	return CFG
}
