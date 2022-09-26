package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"task_axxon/internal/configs"
)

type ServiceImpl struct {
	cli *http.Client
	cfg *configs.Configs
}

func NewService(cli *http.Client,cfg *configs.Configs) Service {
	return &ServiceImpl{cli: cli, cfg: cfg}
}

func (s ServiceImpl) MakeRequest(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {

	//getting url from map
	urlStr, urlOk := req[s.cfg.JsonKeys.Url].(string)
	if !urlOk {
		log.Errorf("missing url key in map")
		return nil, echo.NewHTTPError(http.StatusBadRequest, "url is missing")
	}

	//getting method from map
	method, methodOk := req[s.cfg.JsonKeys.Method].(string)
	if !methodOk {
		log.Errorf("missing method key in map")
		return nil, echo.NewHTTPError(http.StatusBadRequest, "url is missing")
	}

	//making all cap letters
	method = strings.ToUpper(method)

	//creating url struct for make a request
	reqUrl, reqUrlErr := url.Parse(urlStr)
	if reqUrlErr != nil {
		log.Errorf("failed to parse url: %v", reqUrlErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse url: %v", reqUrlErr))
	}

	//getting queries from map
	queries, queriesOk := req[s.cfg.JsonKeys.Queries].(map[string]interface{})
	if queriesOk {
		//creating query struct from url struct above
		query := reqUrl.Query()

		for k, v := range queries {
			qVal, qValOk := v.(string)
			if !qValOk {
				continue
			}
			//adding key and value of queries
			query.Add(k, qVal)
		}
		//mapping encoded queries to url struct
		reqUrl.RawQuery = query.Encode()
	}

	//getting body from map
	body, bodyOk := req[s.cfg.JsonKeys.Body].(map[string]interface{})
	if !bodyOk {
		log.Warnf("missing body key in map")
		//making request without body
		return MakeHTTPRequest(req, method, reqUrl.String(), s.cli, nil, s.cfg.JsonKeys.Headers)
	}

	//if we got body, marshalling from map to bytes
	jsonBody, jbErr := json.Marshal(body)
	if jbErr != nil {
		log.Errorf("invalid body struct: %v", jbErr)
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid body struct: %v", jbErr))
	}

	//making request with body
	return MakeHTTPRequest(req, method, reqUrl.String(), s.cli, bytes.NewBuffer(jsonBody), s.cfg.JsonKeys.Headers)

}

func MakeHTTPRequest(req map[string]interface{}, method, reqUrl string, cli *http.Client, body *bytes.Buffer, headersKey string) (map[string]interface{}, error) {

	//creating newRequest struct
	httpRequest, hReqErr := http.NewRequest(method, reqUrl, nil)
	if hReqErr != nil {
		log.Errorf("failed to create new request: %v", hReqErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create new request: %v", hReqErr))
	}
	if body != nil{
		httpRequest, hReqErr = http.NewRequest(method, reqUrl, body)
		if hReqErr != nil {
			log.Errorf("failed to create new request: %v", hReqErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create new request: %v", hReqErr))
		}
	}

	//getting headers from map
	headers, headersOk := req[headersKey].(map[string]interface{})
	if headersOk {
		//if headers are exists we adding it to request struct
		for k, v := range headers {
			hVal, hValOk := v.(string)
			if !hValOk {
				continue
			}

			httpRequest.Header.Add(k, hVal)
		}
	}
	//making request by http client
	cliResp, cliReqErr := cli.Do(httpRequest)
	if cliReqErr != nil {
		log.Errorf("failed to do request by http client: %v", cliReqErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to do request by http client: %v", cliReqErr))
	}

	//reading body from response
	respBody, respBodyErr := io.ReadAll(cliResp.Body)
	if respBodyErr != nil {
		log.Errorf("failed to read response body: %v", respBodyErr)
		return nil, echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("failed to read response body: %v", respBodyErr))
	}

	//deferred read close function
	defer cliResp.Body.Close()

	//creating map for response
	res := make(map[string]interface{})

	res["status"] = cliResp.StatusCode

	res["length"] = cliResp.ContentLength

	res["headers"] = cliResp.Header

	//map for body
	var resBodyMap map[string]interface{}

	//trying unmarshall body to map
	resBodyUnmErr := json.Unmarshal(respBody, &resBodyMap)
	if resBodyUnmErr != nil {
		log.Warnf("failed to unmarshall response body: %v", resBodyUnmErr)
		//if we could not unmarshall just put it as string
		res["body"] = string(respBody)

		return res, nil
	}

	res["body"] = resBodyMap

	return res, nil
}
