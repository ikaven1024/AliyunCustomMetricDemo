package cms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

const (
	defaultEndpoint = "https://metrichub-cms-cn-hangzhou.aliyuncs.com"
)

type Response struct {
	RequestId string
	Code      string
	Message   string
}

// Client ...
type Client struct {
	Endpoint        string // IP or hostname of SLS endpoint
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
	SourceIP        string
	UserAgent       string // default defaultLogUserAgent
	httpClient      *http.Client
	accessKeyLock sync.RWMutex
}

func NewClient(accessKeyID string, accessKeySecret string) *Client {
	return &Client{
		Endpoint: defaultEndpoint,
		AccessKeyID: accessKeyID,
		AccessKeySecret: accessKeySecret,

		SourceIP:   "192.168.0.112",
		httpClient: http.DefaultClient,
	}
}

func (c *Client) PutCustomMetrics(metrics []CustomMetric) (*Response, error) {
	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}
	uri := "/metric/custom/upload"
	return c.request(http.MethodPost, uri, body)
}

func (c *Client) request(method, uri string, body []byte) (*Response, error) {
	//body = []byte(`[{"dimensions":{"ip":"127.0.0.1","key":"value"},"groupId":102,"metricName":"testMetric","period":15,"time":"20190327T004634.241+0800","type":1,"values":{"LastValue":100}}]`)

	headers := c.buildHeaders()
	if len(body) > 0 {
		headers["Content-MD5"] = md5Str(body)
	}
	headers["Content-Length"] = fmt.Sprintf("%d", len(body))

	err := c.signature(method, uri, headers)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, c.Endpoint + uri, bytes.NewBuffer(body))
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	if err != nil {
		return nil, err
	}
	httpResp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(httpResp.Body)

	resp := &Response{}
	err = json.Unmarshal(bodyBytes, resp)
	if err != nil {
		return nil, fmt.Errorf("%s %s错误： %d %s", method, uri, httpResp.StatusCode, string(bodyBytes))
	}
	return resp, nil
}

func (c *Client) buildHeaders() map[string]string {
	headers := map[string]string {
		"User-Agent":        "cms-go-sdk-v-1.0",
		"Content-Length":    "0",
		"Content-Type":      "application/json",
		"Date":              nowRFC1123(),
		"Host":              "metrichub-cms-cn-hangzhou.aliyuncs.com",
		"x-cms-api-version": "1.0",
		"x-cms-signature":   "hmac-sha1",
		"x-cms-ip":          c.SourceIP,
	}
	if len(c.SecurityToken) > 0 {
		headers["x-cms-caller-type"] = "token"
		headers["x-cms-security-token"] = c.SecurityToken
	}
	return headers
}

func (c *Client) signature(method, uri string, headers map[string]string) error {
	var contentMD5, contentType, date, canoHeaders, canoResource string
	var slsHeaderKeys sort.StringSlice

	if val, ok := headers["Content-MD5"]; ok {
		contentMD5 = val
	}

	if val, ok := headers["Content-Type"]; ok {
		contentType = val
	}

	date, ok := headers["Date"]
	if !ok {
		return fmt.Errorf("can't find 'Date' header")
	}

	// Calc CanonicalizedSLSHeaders
	slsHeaders := make(map[string]string, len(headers))
	for k, v := range headers {
		l := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(l, "x-cms-") || strings.HasPrefix(l, "x-acs-") {
			slsHeaders[l] = strings.TrimSpace(v)
			slsHeaderKeys = append(slsHeaderKeys, l)
		}
	}

	sort.Sort(slsHeaderKeys)
	for i, k := range slsHeaderKeys {
		canoHeaders += k + ":" + slsHeaders[k]
		if i+1 < len(slsHeaderKeys) {
			canoHeaders += "\n"
		}
	}

	// Calc CanonicalizedResource
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	canoResource += u.EscapedPath()
	if u.RawQuery != "" {
		var keys sort.StringSlice

		vals := u.Query()
		for k := range vals {
			keys = append(keys, k)
		}

		sort.Sort(keys)
		canoResource += "?"
		for i, k := range keys {
			if i > 0 {
				canoResource += "&"
			}

			for _, v := range vals[k] {
				canoResource += k + "=" + v
			}
		}
	}

	signStr := method + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		canoHeaders + "\n" +
		canoResource

	// Signature = base64(hmac-sha1(UTF8-Encoding-Of(SignString)，AccessKeySecret))
	mac := hmac.New(sha1.New, []byte(c.AccessKeySecret))
	_, err = mac.Write([]byte(signStr))
	if err != nil {
		return err
	}
	//digest := base32.StdEncoding.EncodeToString(mac.Sum(nil))
	digest := base16(mac.Sum(nil))
	headers["Authorization"] = c.AccessKeyID + ":" + digest
	return nil
}
