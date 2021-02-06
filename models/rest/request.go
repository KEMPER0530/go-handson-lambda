package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

func FetchRequest(url string, m string) (int, []byte) {
	cli, err := newClient()
	if err != nil {
		log.Fatalln("Fail to make http.Client", err)
	}

	req, err := newRequest(m, url, nil)
	if err != nil {
		log.Fatalln("Fail to make http.Request", err)
	}

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)

	resp, err := do(cli, req, nil)
	if err != nil {
		log.Fatal("Fail on HTTP request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return resp.StatusCode, nil
	}

	// レスポンスボディをすべて読み出す
	body, _ := ioutil.ReadAll(resp.Body)

	// ステータスコード
	fmt.Println("status:", resp.StatusCode)
	// ヘッダーを取得
	fmt.Println("header:", resp.Header.Get("Content-Type"))
	// resp.Bodyの大きさ len(body) と同じ
	fmt.Println("ContentLength:", resp.ContentLength)
	// リクエストURL
	fmt.Println("Request:", resp.Request.URL.String())

	return resp.StatusCode, body
}

func newClient() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}
	return client, nil
}

func newRequest(method, path string, values url.Values) (*http.Request, error) {
	body := strings.NewReader(values.Encode())
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func do(client *http.Client, req *http.Request, v interface{}) (*http.Response, error) {

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			// give *bytes.Buffer to get raw bytes instead of json decoded string
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)

			// ignore the error caused by an empty response
			if err == io.EOF {
				err = nil
			}
		}
	}

	return resp, nil
}
