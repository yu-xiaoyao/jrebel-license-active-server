package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
)

func main() {

	content := "<PingResponse><message></message><responseCode>OK</responseCode><salt>ABCD</salt></PingResponse>"

	signatures, err := signWithMd5([]byte(content))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(signatures))

	signature, err := signWithSha1([]byte{0, 0, 0, 0})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(signature)

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
