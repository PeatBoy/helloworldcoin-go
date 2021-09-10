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
func ServiceUnauthorized(rw http.ResponseWriter) {
	fail(rw, SERVICE_UNAUTHORIZED)
}

const REQUEST_PARAM_ILLEGAL = "request_param_illegal"
const SERVICE_UNAVAILABLE = "service_unavailable"
const SERVICE_UNAUTHORIZED = "service_unauthorized"

const NOT_FOUND_TRANSACTION = "not_found_transaction"
const NOT_FOUND_BLOCK = "not_found_block"
const NOT_FOUND_UNCONFIRMED_TRANSACTIONS = "not_found_unconfirmed_transactions"
