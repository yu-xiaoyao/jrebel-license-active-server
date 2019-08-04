package main

import "encoding/json"

var jRebelLeases JRebelLeasesStruct
var jrebelLeases1 JrebelLeases1Struct
var jrebelValidate JrebelValidateStruct

func init() {
	_ = json.Unmarshal([]byte(jrebelLeasesJson), &jRebelLeases)
	_ = json.Unmarshal([]byte(jrebelLeases1Json), &jrebelLeases1)
	_ = json.Unmarshal([]byte(jrebelValidateJson), &jrebelValidate)
}

//language=JSON
const jrebelLeasesJson = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "id": 1,
    "licenseType": 1,
    "evaluationLicense": false,
    "signature": "OJE9wGg2xncSb+VgnYT+9HGCFaLOk28tneMFhCbpVMKoC/Iq4LuaDKPirBjG4o394/UjCDGgTBpIrzcXNPdVxVr8PnQzpy7ZSToGO8wv/KIWZT9/ba7bDbA8/RZ4B37YkCeXhjaixpmoyz/CIZMnei4q7oWR7DYUOlOcEWDQhiY=",
    "serverRandomness": "H2ulzLlh7E0=",
    "seatPoolType": "standalone",
    "statusCode": "SUCCESS",
    "offline": false,
    "validFrom": null,
    "validUntil": null,
    "company": "Administrator",
    "orderId": "",
    "zeroIds": [
        
    ],
    "licenseValidFrom": 1490544001000,
    "licenseValidUntil": 1691839999000
}`

type JRebelLeasesStruct struct {
	ServerVersion         string        `json:"serverVersion"`
	ServerProtocolVersion string        `json:"serverProtocolVersion"`
	ServerGUID            string        `json:"serverGuid"`
	GroupType             string        `json:"groupType"`
	ID                    int           `json:"id"`
	LicenseType           int           `json:"licenseType"`
	EvaluationLicense     bool          `json:"evaluationLicense"`
	Signature             string        `json:"signature"`
	ServerRandomness      string        `json:"serverRandomness"`
	SeatPoolType          string        `json:"seatPoolType"`
	StatusCode            string        `json:"statusCode"`
	Offline               bool          `json:"offline"`
	ValidFrom             int64         `json:"validFrom"`
	ValidUntil            int64         `json:"validUntil"`
	Company               string        `json:"company"`
	OrderID               string        `json:"orderId"`
	ZeroIds               []interface{} `json:"zeroIds"`
	LicenseValidFrom      int64         `json:"licenseValidFrom"`
	LicenseValidUntil     int64         `json:"licenseValidUntil"`
}

//language=JSON
const jrebelLeases1Json = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "statusCode": "SUCCESS",
    "msg": null,
    "statusMessage": null
}
`

type JrebelLeases1Struct struct {
	ServerVersion         string      `json:"serverVersion"`
	ServerProtocolVersion string      `json:"serverProtocolVersion"`
	ServerGUID            string      `json:"serverGuid"`
	GroupType             string      `json:"groupType"`
	StatusCode            string      `json:"statusCode"`
	Company               string      `json:"company"`
	Msg                   interface{} `json:"msg"`
	StatusMessage         interface{} `json:"statusMessage"`
}

//language=JSON
const jrebelValidateJson = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "statusCode": "SUCCESS",
    "company": "Administrator",
    "canGetLease": true,
    "licenseType": 1,
    "evaluationLicense": false,
    "seatPoolType": "standalone"
}
`

type JrebelValidateStruct struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGUID            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Company               string `json:"company"`
	CanGetLease           bool   `json:"canGetLease"`
	LicenseType           int    `json:"licenseType"`
	EvaluationLicense     bool   `json:"evaluationLicense"`
	SeatPoolType          string `json:"seatPoolType"`
}
