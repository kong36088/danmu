package danmu

import(
	log "github.com/alecthomas/log4go"
)

func InitLog() error{
	logPath := Conf.GetConfig("sys","log_conf")
	log.LoadConfiguration(logPath)
	return OK
}

func CloseLog(){
	log.Close()
}
