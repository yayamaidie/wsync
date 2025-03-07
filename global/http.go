package global

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func HttpRequest(_method string, _url string, _body interface{}) (_statuscode int, _respbytes []byte, _err error) {
	var req *http.Request

	reqbodybytes, err := json.Marshal(_body)
	if err != nil {
		return -1, nil, err
	}
	bytereader := bytes.NewReader(reqbodybytes)

	switch _method {
	case "GET":
		req, err = http.NewRequest("GET", _url, nil)
		if err != nil {
			return -1, nil, err
		}
	case "POST":
		if _body == nil {
			req, err = http.NewRequest("POST", _url, nil)
			if err != nil {
				return -1, nil, err
			}
		} else {
			req, err = http.NewRequest("POST", _url, bytereader)
			if err != nil {
				return -1, nil, err
			}
		}
	case "PUT":
		req, err = http.NewRequest("PUT", _url, bytereader)
		if err != nil {
			return -1, nil, err
		}
	case "DELETE":
		req, err = http.NewRequest("DELETE", _url, nil)
		if err != nil {
			return -1, nil, err
		}
	}

	block, _ := pem.Decode(CertBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return -1, nil, fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return -1, nil, err
	}
	cert_pool := x509.NewCertPool()
	cert_pool.AddCert(cert)

	httpclient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cert_pool,
				// InsecureSkipVerify: true,
			},
		},
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}

	return resp.StatusCode, respbytes, nil
}

func ServerIsHttp(rawurl string) (bool, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return false, err
	}

	url_string := fmt.Sprintf("http://%s", net.JoinHostPort(url.Hostname(), url.Port()))
	req, err := http.NewRequest("GET", url_string, nil)
	if err != nil {
		return false, err
	}

	hc := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := hc.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 && strings.Contains(string(respbytes), "Client sent an HTTP request to an HTTPS server") {
		return false, nil
	}
	return true, nil
}
