package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func loggingRequest(tag string, r *http.Request) {
	fmt.Printf("%s --- %s\n", time.Now(), tag)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest("indexHandler", r)
	host := "http://" + r.Host

	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	html := `<h1>Hello,This is a Jrebel & JetBrains License Server!</h1>
<p>License Server started at %s
<p>JRebel 7.1 and earlier version Activation address was: <span style='color:red'>%s/{tokenname}</span>, with any email."
<p>JRebel 2018.1 and later version Activation address was: %s/{guid}(eg:<span style='color:red'> %s/%s </span>), with any email.`
	_, _ = fmt.Fprintf(w, html, host, host, host, host, newUUIDV4String())
}

func jrebelLeasesHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest("jrebelLeasesHandler", r)

	w.Header().Set("content-type", "application/json; charset=utf-8")

	parameter, err := getHttpBodyParameter(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}

	clientRandomness := parameter.Get("randomness")
	username := parameter.Get("username")
	guid := parameter.Get("guid")
	if clientRandomness == "" || username == "" || guid == "" {
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
	var responseBody = jRebelLeases
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

		responseBody.Offline = offline
		responseBody.ValidFrom, _ = strconv.ParseInt(validFrom, 10, 64)
		responseBody.ValidUntil, _ = strconv.ParseInt(validUntil, 10, 64)
	}
	serverRandomness := newServerRandomness()
	signature := toLeaseCreateJson(clientRandomness, serverRandomness, guid, offline, validFrom, validUntil)

	responseBody.ServerRandomness = serverRandomness
	responseBody.Signature = signature
	responseBody.Company = username

	response(w, &responseBody)
}

func jrebelLeases1Handler(w http.ResponseWriter, r *http.Request) {
	loggingRequest("jrebelLeases1Handler", r)

	w.Header().Set("content-type", "application/json; charset=utf-8")
	parameter, err := getHttpBodyParameter(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}
	username := parameter.Get("username")

	var responseBody = jrebelLeases1
	if username != "" {
		responseBody.Company = username
	}

	response(w, &responseBody)
}

func jrebelValidateHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest("jrebelValidateHandler", r)

	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "%s\n", jrebelValidateJson)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest("pingHandler", r)

	w.Header().Add("content-type", "text/html; charset=utf-8")
	parameter, err := getHttpBodyParameter(r)
	if err != nil {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}
	salt := parameter.Get("salt")
	if salt == "" {
		w.WriteHeader(403)
		_, _ = fmt.Fprint(w)
	} else {
		xmlContent := "<PingResponse><message></message><responseCode>OK</responseCode><salt>" + salt + "</salt></PingResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(403)
			_, _ = fmt.Fprintf(w, "%s\n", err)
		} else {
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "%s\n", body)
		}
	}
}

func obtainTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json; charset=utf-8")

	parameter, err := getHttpBodyParameter(r)
	if err != nil {
		responseError(w, err, 403)
		return
	}
	salt := parameter.Get("salt")
	username := parameter.Get("userName")
	prolongationPeriod := "607875500"
	if salt == "" || username == "" {
		w.WriteHeader(403)
		_, _ = fmt.Fprintln(w)
	} else {
		w.WriteHeader(200)
		xmlContent := "<ObtainTicketResponse><message></message><prolongationPeriod>" + prolongationPeriod + "</prolongationPeriod><responseCode>OK</responseCode><salt>" + salt + "</salt><ticketId>1</ticketId><ticketProperties>licensee=" + username + "\tlicenseType=0\t</ticketProperties></ObtainTicketResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(403)
			_, _ = fmt.Fprintf(w, "%s\n", err)
		} else {
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "%s\n", body)
		}
	}

}
func releaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html; charset=utf-8")
	parameter, err := getHttpBodyParameter(r)
	if err != nil {
		responseError(w, err, 403)
		return
	}
	salt := parameter.Get("salt")
	if salt == "" {
		w.WriteHeader(403)
		_, _ = fmt.Fprintln(w)
	} else {
		xmlContent := "<ReleaseTicketResponse><message></message><responseCode>OK</responseCode><salt>" + salt + "</salt></ReleaseTicketResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(403)
			_, _ = fmt.Fprintf(w, "%s\n", err)
		} else {
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "%s\n", body)
		}
	}

}

func responseError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(403)
	_, _ = fmt.Fprintf(w, "%s\n", err)
}

func response(w http.ResponseWriter, resp interface{}) {
	bodyData, err := json.Marshal(&resp)
	if err != nil {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "%s\n", string(bodyData))
}

func getHttpBodyParameter(r *http.Request) (params url.Values, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
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
	return ps.Query(), err
}
