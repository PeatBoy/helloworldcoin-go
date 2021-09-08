package controller

import (
	"helloworld-blockchain-go/util/JsonUtil"
	"io"
	"io/ioutil"
	"net/http"
)

func GetRequest(req *http.Request, object interface{}) interface{} {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), object)
	return request
}

func success(rw http.ResponseWriter, response interface{}) {
	s := "{\"status\":\"success\",\"message\":null,\"data\":" + JsonUtil.ToString(response) + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func fail(rw http.ResponseWriter, message string) {
	s := "{\"status\":\"fail\",\"message\":\"" + message + "\",\"data\":null" + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func requestParamIllegal(rw http.ResponseWriter) {
	fail(rw, REQUEST_PARAM_ILLEGAL)
}
func serviceUnavailable(rw http.ResponseWriter) {
	fail(rw, SERVICE_UNAVAILABLE)
}

const REQUEST_PARAM_ILLEGAL = "REQUEST_PARAM_ILLEGAL"
const SERVICE_UNAVAILABLE = "SERVICE_UNAVAILABLE"

const NOT_FOUND_TRANSACTION = "NOT_FOUND_TRANSACTION"
const NOT_FOUNT_BLOCK = "NOT_FOUNT_BLOCK"
const NOT_FOUNT_UNCONFIRMED_TRANSACTIONS = "NOT_FOUNT_UNCONFIRMED_TRANSACTIONS"
