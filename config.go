package danmu

import (
	"flag"
	"runtime"
	"github.com/larspensjo/config"
	"log"
)

var (
	configFile = flag.String("config", "config/config.ini", "General configuration file")
)

//topic list
var CFG map[string]string

func ReadConfig() map[string]string {
	if CFG != nil {
		return CFG
	}

	CFG = make(map[string]string)

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
		for _, v := range options {
			option, err := cfg.String(section, v)
			if err != nil {
				log.Printf("Read file %s option %s failed, %s\n", *configFile, v, err)
				continue
			}
			CFG[v] = option
		}
	}

	return CFG
}

func GetConfig(option string) string {
	if CFG == nil {
		ReadConfig()
	}
	return CFG[option]
}

func GetAllConfig() map[string]string {
	if CFG == nil {
		ReadConfig()
	}
	return CFG
}
