package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	serverPort := 12345
	if len(os.Args) > 1 {
		for k, v := range os.Args {
			if k == 0 {
				continue
			}
			arg := strings.TrimSpace(v)
			hasPortL := strings.HasPrefix(arg, "--port=")
			hasPortS := strings.HasPrefix(arg, "-p=")
			if hasPortL {
				i, err := strconv.ParseInt(strings.ReplaceAll(arg, "--port=", ""), 10, 32)
				if err == nil {
					serverPort = int(i)
					break
				}
			}
			if hasPortS {
				i, err := strconv.ParseInt(strings.ReplaceAll(arg, "-p=", ""), 10, 32)
				if err == nil {
					serverPort = int(i)
					break
				}
			}
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/jrebel/leases", jrebelLeasesHandler)
	http.HandleFunc("/jrebel/leases/1", jrebelLeases1Handler)
	http.HandleFunc("/agent/leases", jrebelLeasesHandler)
	http.HandleFunc("/agent/leases/1", jrebelLeases1Handler)
	http.HandleFunc("/jrebel/validate-connection", jrebelValidateHandler)
	http.HandleFunc("/rpc/ping.action", pingHandler)
	http.HandleFunc("/rpc/obtainTicket.action", obtainTicketHandler)
	http.HandleFunc("/rpc/releaseTicket.action", releaseTicketHandler)

	fmt.Printf("start server with port = %d\n", serverPort)

	_ = http.ListenAndServe(":"+strconv.Itoa(serverPort), nil)
}
