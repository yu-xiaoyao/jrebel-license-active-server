// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	OfflineDays    int
	LogLevel       int64
	ExportSchema   string
	ExportHost     string
	NewIndex       bool
	JrebelWorkMode int
}

var config = &Config{
	Port:           12345,
	OfflineDays:    30,
	LogLevel:       Info,
	ExportSchema:   "http",
	ExportHost:     "",   // default is request ip
	NewIndex:       true, // use new index page
	JrebelWorkMode: 0,    // 0: auto, 1: force offline mode, 2: force online mode, 3: oldGuid
}

var logger = NewLogger(os.Stdout, Info, log.Ldate|log.Ltime)

func init() {
	portPtr := flag.Int("port", config.Port, "Server port, value range 1-65535")
	logLevelPtr := flag.Int64("logLevel", config.LogLevel, "Log level, value range 1-4")
	exportSchemaPtr := flag.String("exportSchema", config.ExportSchema, "Index page export schema (http or https)")
	exportHostPtr := flag.String("exportHost", "", "Index page export host, ip or domain")
	newIndexPtr := flag.Bool("newIndex", config.NewIndex, "Use New Index Page (true or false)")
	jrebelWorkMode := flag.Int("jrebelWorkMode", config.JrebelWorkMode, "Jrebel Wrok mode. 0: auto, 1: force offline mode, 2: force online mode, 3: oldGuid")

	flag.Parse()

	config.Port = *portPtr
	config.LogLevel = *logLevelPtr
	config.ExportSchema = *exportSchemaPtr
	config.ExportHost = *exportHostPtr
	config.NewIndex = *newIndexPtr
	config.JrebelWorkMode = *jrebelWorkMode

	logger.SetLevel(config.LogLevel)
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
