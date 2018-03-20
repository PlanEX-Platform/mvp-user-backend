package logenv

import (
	"os"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLog() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// var logPath = "./planex-" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	logDir := viper.GetString("logs.dir")
	logPath := logDir + "/planex-dbg.log"
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second))

	if err != nil {
		log.Error(err)
	}

	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.ErrorLevel: writer,
		},
		&log.TextFormatter{}))
	// &log.JSONFormatter{}))
}
