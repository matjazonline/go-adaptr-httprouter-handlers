package adaptr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetApiKeyFromReq(r *http.Request) string {
	apiKeyStr := r.URL.Query().Get("apiKey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.URL.Query().Get("apikey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.FormValue("apiKey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	apiKeyStr = r.FormValue("apikey")
	if apiKeyStr != "" {
		return apiKeyStr
	}

	return ""
}

func GetParamFromReqString(r *http.Request, paramName string) string {
	qStr := r.URL.Query().Get(paramName)
	if qStr != "" {
		return qStr
	}
	return ""
}

func GetTokenFromReq(r *http.Request, tknParameterName string, requestJsonStructCtxKey interface{}) (string, error) {
	if tknParameterName != "" {
		var tknParValue string
		if requestJsonStructCtxKey == nil {
			requestJsonStructCtxKey = CtxRequestJsonStructKey
		}
		ctxJsonStruct := GetCtxValue(r, requestJsonStructCtxKey).(map[string]interface{})
		if ctxJsonStruct != nil {
			if v, ok := ctxJsonStruct[tknParameterName]; ok {
				tknParValue = v.(string)
			}
		}

		if tknParValue == "" {
			tknParValue = r.FormValue(tknParameterName)
			if tknParValue == "" {
				return "", fmt.Errorf("no token value in parameter=%v", tknParameterName)
			}
		}

		return tknParValue, nil
	}

	authHeaderVal := r.Header.Get("Authorization")
	if authHeaderVal == "" {
		return "", errors.New("No Authorization header value")
	}

	bearerStr := "Bearer"
	if last := strings.LastIndex(authHeaderVal, bearerStr); last > -1 {
		tknValue := strings.TrimSpace(authHeaderVal[last+len(bearerStr):])
		if tknValue != "" {
			return tknValue, nil
		}
	}

	return "", errors.New("Authorization header parse failed")
}

func JsonOut(w http.ResponseWriter, jsonOutPointer interface{}) {
	res, _ := json.Marshal(jsonOutPointer)
	fmt.Fprint(w, string(res))
}

func GetHttprouterQueryParameterValue(r *http.Request, name string) (string, error) {
	params, ok := GetCtxValue(r, CtxHttpRouterParamsKey).(httprouter.Params)
	if ok {
		return params.ByName(name), nil
	}
	return "", errors.New("no parameters in context")
}

func getBodyValue(r *http.Request) (*http.Request, error) {
	bodyVal := GetCtxValue(r, CtxRequestBodyByteArrKey)
	if bodyVal != nil {
		return nil, nil
	}

	bodyValRead, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return SetCtxValue(r, CtxRequestBodyByteArrKey, bodyValRead), nil
}
