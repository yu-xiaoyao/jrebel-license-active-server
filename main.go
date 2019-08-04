package main

import (
	"net/http"
)

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

	_ = http.ListenAndServe(":12345", nil)
}
