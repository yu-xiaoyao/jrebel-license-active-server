package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	signature, err := testSign([]byte{0, 0, 0, 0})
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	//port := 1000
	//_, _ = fmt.Fprintf(w, html)
}

func jrebelLeasesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(len(r.PostForm))
	fmt.Println(len(r.Form))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		parameter := toHttpBodyParameter(body)
		randomness := parameter.Get("randomness")
		username := parameter.Get("username")
		guid := parameter.Get("guid")
		if randomness == "" || username == "" || guid == "" {
			w.WriteHeader(403)
			_, _ = fmt.Fprint(w)
			return
		}
		offline, err := strconv.ParseBool(parameter.Get("offline"))
		if err != nil {
			offline = false
		}

		validFrom := "null"
		validUntil := "null"
		if offline {
			clientTime := parameter.Get("clientTime")
			_ = parameter.Get("offlineDays")

			startTimeInt, err := strconv.ParseInt(clientTime, 10, 64)
			if err != nil {
				startTimeInt = int64(time.Now().Second()) * 1000
			}
			// 过期时间
			expTime := int64(180 * 24 * 60 * 60 * 100)
			validFrom = clientTime
			validUntil = strconv.FormatInt(startTimeInt+expTime, 10)
		}

		fmt.Println(validFrom)
		fmt.Println(validUntil)

	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	fmt.Fprintf(w, "Hello there!\n")
}

func jrebelLeases1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("jrebelLeases1Handler")

	randomness := r.Form.Get("randomness")
	fmt.Println(randomness)

	fmt.Fprintf(w, "Hello there!\n")
}

func jrebelValidateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there!\n")
}
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there!\n")
}

func obtainTicketHandler(w http.ResponseWriter, r *http.Request) {

}
func releaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there!\n")
}

func toHttpBodyParameter(body []byte) url.Values {
	s := string(body)
	ps := url.URL{
		Scheme:     "",
		Opaque:     "",
		User:       nil,
		Host:       "",
		Path:       "",
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   s,
		Fragment:   "",
	}
	fmt.Println(s)
	return ps.Query()
}
