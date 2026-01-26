// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"flag"
	"net/http"
	"strconv"
)

type Config struct {
	Port             int
	OfflineDays      int64
	IgnoreOfflineDay bool
	LogLevel         Level
	LogFile          bool
	LogPath          string
	ExportSchema     string
	ExportHost       string
	NewIndex         bool
	JrebelWorkMode   int
}

var config = &Config{
	Port:             12345,
	IgnoreOfflineDay: false,
	OfflineDays:      30, // max 180 > 180 will cause invalid
	LogLevel:         Info,
	LogFile:          false,
	LogPath:          "./logs",
	ExportSchema:     "http",
	ExportHost:       "",   // default is request ip
	NewIndex:         true, // use new index page
	JrebelWorkMode:   0,    // 0: auto, 1: force offline mode, 2: force online mode, 3: oldGuid
}

// var logger = NewLogger(os.Stdout, Info, log.Ldate|log.Ltime)
var logger ILogger

func init() {
	portPtr := flag.Int("port", config.Port, "Server port, value range 1-65535")
	ignoreOfflineDay := flag.Bool("ignoreOfflineDay", config.IgnoreOfflineDay, "Ignore plugin offline day parameter, if true force return offlineDays parameter. default: false")
	offlineDays := flag.Int64("offlineDays", config.OfflineDays, "Custom return offline days parameter, Recommended not to exceed 180 days. default: 30")
	//logLevelPtr := flag.Int64("logLevel", config.LogLevel, "Log level, value range 1-4")
	var logLevel = Info
	flag.Var(&logLevel, "level", "Log level, value range 0-5")

	logFile := flag.Bool("logFile", config.LogFile, "Log to File")
	logPath := flag.String("logPath", config.LogPath, "Save log file path")

	exportSchemaPtr := flag.String("exportSchema", config.ExportSchema, "Index page export schema (http or https)")
	exportHostPtr := flag.String("exportHost", "", "Index page export host, ip or domain")
	newIndexPtr := flag.Bool("newIndex", config.NewIndex, "Use New Index Page (true or false)")
	jrebelWorkMode := flag.Int("jrebelWorkMode", config.JrebelWorkMode, "Jrebel Work mode. 0: auto, 1: force offline mode, 2: force online mode, 3: oldGuid")

	flag.Parse()

	config.Port = *portPtr
	config.IgnoreOfflineDay = *ignoreOfflineDay
	config.OfflineDays = *offlineDays
	config.LogFile = *logFile
	config.LogPath = *logPath
	config.ExportSchema = *exportSchemaPtr
	config.ExportHost = *exportHostPtr
	config.NewIndex = *newIndexPtr
	config.JrebelWorkMode = *jrebelWorkMode

	config.LogLevel = logLevel

	if config.LogFile {
		logger = NewFileLogger(*logPath, 100, 7, logLevel, 0)
	} else {
		logger = NewConsoleLogger(logLevel, 0)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/jrebel/leases", jrebelLeasesHandler)
	http.HandleFunc("/jrebel/leases/1", jrebelLeases1Handler)
	http.HandleFunc("/agent/leases", jrebelLeasesHandler)
	http.HandleFunc("/agent/leases/1", jrebelLeases1Handler)
	http.HandleFunc("/jrebel/validate-connection", jrebelValidateHandler)
	http.HandleFunc("/rpc/ping.action", pingHandler)
	http.HandleFunc("/rpc/obtainTicket.action", obtainTicketHandler)
	http.HandleFunc("/rpc/releaseTicket.action", releaseTicketHandler)

	logger.Infof("Start server with port = %d, schema = %s\n", config.Port, config.ExportSchema)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		logger.Errorf("Start server failed. cause: %v\n", err)
	}
}
