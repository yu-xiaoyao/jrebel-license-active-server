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
	OfflineDefault bool
	OfflineDays    int
	LogLevel       int64
	Schema         string // 新增 schema 字段
}

var config = &Config{
	Port:           12345,
	OfflineDefault: true,
	OfflineDays:    30,
	LogLevel:       Info,
	Schema:         "http", // 默认协议为 http,https
}

var logger = NewLogger(os.Stdout, Info, log.Ldate|log.Ltime)

func init() {
	portPtr := flag.Int("port", config.Port, "Server port, value range 1-65535")
	logLevelPtr := flag.Int64("logLevel", config.LogLevel, "Log level, value range 1-4")
	schemaPtr := flag.String("schema", config.Schema, "Protocol schema (http or https)")

	flag.Parse()

	config.Port = *portPtr
	config.LogLevel = *logLevelPtr
	config.Schema = *schemaPtr

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

	logger.Infof("Start server with port = %d, schema = %s\n", config.Port, config.Schema)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		logger.Errorf("Start server failed. cause: %v\n", err)
	}
}
