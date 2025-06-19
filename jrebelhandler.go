// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var indexTemplate *template.Template

func init() {
	indexTemplate, _ = template.New("index").Parse(indexTemplateHtml)
}

func loggingRequest(r *http.Request) {
	query := r.URL.RawQuery
	if query != "" {
		query = "?" + query
	}

	logger.Infof("--> %s %s%s. [%s] [%s]\n", r.Method, r.URL.Path, query, r.RemoteAddr, r.UserAgent())

	// debug info
	//contentType := r.Header.Get("Content-Type")
	//logger.Debugf("Content-Type: %s\n", contentType)
	//logger.Debugf("Host: %s\n", r.Host)
	//logger.Debugf("RemoteAddr: %s\n", r.RemoteAddr)
	//logger.Debugf("User-Agent: %s\n", r.UserAgent())

	if r.Method == "POST" {
		if logger.IsDebug() {
			body, _ := io.ReadAll(r.Body)
			logger.Debugf("--> Request Body: %s", body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
}

func jrebelLeasesHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest(r)

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

	// active by product JRebel = offline, XRebel = online
	var offline bool
	product := parameter.Get("product")
	if product == "XRebel" {
		offline = false
	} else {
		offline, err = strconv.ParseBool(parameter.Get("offline"))
		if err != nil {
			offline = true
		}
		oldGuid := parameter.Get("oldGuid")
		if oldGuid != "" {
			offline = true
		}
	}

	validFrom := ""
	validUntil := ""
	var responseBody = jRebelLeases
	if offline {
		clientTime := parameter.Get("clientTime")
		offlineDays := parameter.Get("offlineDays")

		startTimeInt, err := strconv.ParseInt(clientTime, 10, 64)
		if err != nil {
			startTimeInt = int64(time.Now().Second()) * 1000
		}

		offlineDaysInt, err := strconv.ParseInt(offlineDays, 10, 64)
		if err != nil {
			offlineDaysInt = int64(config.OfflineDays)
		}

		// 过期时间
		expireTime := startTimeInt + (offlineDaysInt * 24 * 60 * 60 * 1000)
		responseBody.Offline = offline
		responseBody.ValidFrom = startTimeInt
		responseBody.ValidUntil = expireTime

		validFrom = clientTime
		validUntil = strconv.FormatInt(expireTime, 10)
	}

	serverRandomness := newServerRandomness()
	signature := toLeaseCreateJson(clientRandomness, serverRandomness, guid, offline, validFrom, validUntil)

	responseBody.ServerRandomness = serverRandomness
	responseBody.Signature = signature
	responseBody.Company = username

	response(w, &responseBody)
}

func jrebelLeases1Handler(w http.ResponseWriter, r *http.Request) {
	loggingRequest(r)

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
	loggingRequest(r)

	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "%s\n", jrebelValidateJson)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest(r)

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
	body, err := io.ReadAll(r.Body)
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
	// fmt.Println(s)
	return ps.Query(), err
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loggingRequest(r)

	var host string
	if len(config.ExportHost) == 0 {
		host = config.ExportSchema + "://" + r.Host
	} else {
		host = config.ExportSchema + "://" + config.ExportHost
	}
	uuid := newUUIDV4String()

	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(200)

	if config.NewIndex {
		data := struct {
			Host string
			UUID string
		}{
			Host: host,
			UUID: uuid,
		}
		err := indexTemplate.Execute(w, data)
		if err == nil {
			return
		} else {
			logger.Warnf("template execute error: %v\n", err)
		}
	}

	// template error, fallback simple
	html := `<h1>Hello,This is a Jrebel License Server!</h1>
			<p>License Server started at %s
			<p>JRebel 7.1 and earlier version Activation address was: <span style='color:red'>%s/{tokenname}</span>, with any email."
			<p>JRebel 2018.1 and later version Activation address was: %s/{guid}(eg:<span style='color:red'> %s/%s </span>), with any email.`
	_, _ = fmt.Fprintf(w, html, host, host, host, host, uuid)

}

const indexTemplateHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>JRebel License Server</title>
    <style>
        :root {
            --primary-color: #4a6bff;
            --secondary-color: #f5f5f5;
            --accent-color: #ff5252;
            --text-color: #333;
            --light-text: #666;
            --border-radius: 8px;
            --box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            --transition: all 0.3s ease;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 900px;
            margin: 40px auto;
            background-color: white;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
            overflow: hidden;
        }

        header {
            background-color: var(--primary-color);
            color: white;
            padding: 20px 30px;
            position: relative;
        }
        .language-btn:hover, .language-btn.active {
            background: rgba(255, 255, 255, 0.4);
        }

        h1 {
            font-size: 28px;
            margin-bottom: 10px;
        }

        .content {
            padding: 30px;
        }

        .info-card {
            background-color: var(--secondary-color);
            border-radius: var(--border-radius);
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }

        .info-card h2 {
            font-size: 20px;
            margin-bottom: 15px;
            color: var(--primary-color);
        }

        .activation-url {
            display: flex;
            align-items: center;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 10px;
            margin: 10px 0;
            position: relative;
        }

        .url-text {
            flex-grow: 1;
            font-family: monospace;
            word-break: break-all;
        }

        .highlight {
            color: var(--accent-color);
            font-weight: bold;
        }

        footer {
            text-align: center;
            padding: 20px;
            color: var(--light-text);
            font-size: 14px;
            border-top: 1px solid #eee;
        }

        footer a {
            color: var(--primary-color);
            text-decoration: none;
        }

        footer a:hover {
            text-decoration: underline;
        }

        @media (max-width: 768px) {
            .container {
                margin: 20px auto;
            }

            header {
                padding: 15px 20px;
            }

            .content {
                padding: 20px;
            }

            .language-switch {
                position: static;
                margin-top: 10px;
                text-align: right;
            }
        }
    </style>
</head>
<body>
<div class="container">
    <header>
        <h1>JRebel License Server</h1>
    </header>

    <div class="content">
        <div class="info-card">
            <h2>Server Information</h2>
            <p>License Server started at:</p>
            <div class="activation-url">
                <span class="url-text">{{.Host}}</span>
            </div>
        </div>

        <div class="info-card">
            <h2>JRebel 7.1 and Earlier Versions</h2>
            <p>Activation address (with any email):</p>
            <div class="activation-url">
                <span class="url-text">{{.Host}}/<span class="highlight">{tokenname}</span></span>
            </div>
        </div>

        <div class="info-card">
            <h2>JRebel 2018.1 and Later Versions</h2>
            <p>Activation address (with any email address):</p>
            <div class="activation-url">
                <span class="url-text">{{.Host}}/<span class="highlight">{{.UUID}}</span></span>
            </div>
        </div>
    </div>

    <footer>
        <p>
            <a href="https://github.com/yu-xiaoyao/jrebel-license-active-server" target="_blank">Developed from 2019
                year</a>
        </p>
    </footer>
</div>

</body>
</html>`
