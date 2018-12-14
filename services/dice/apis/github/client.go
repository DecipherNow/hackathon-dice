package github

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// GetMembers will retrieve all members for an organization from GitHub
func (c *Client) GetMembers(pageNumber int) ([]Member, error) {
	// logger := zerolog.New(os.Stdout).
	// 	With().Timestamp().Str("service", "dice").Logger().
	// 	Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// TODO: iterate all pages building up a response array

	var ret []Member
	uri := viper.GetString("github_api") + "orgs/" + viper.GetString("github_org") + "/members?page=" + strconv.Itoa(pageNumber)
	resp, err := c.doGet(uri, nil)
	if err != nil {
		return ret, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ret, errorFromResponse(resp)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	jsonErr := json.Unmarshal(body, &ret)
	if jsonErr != nil {
		return ret, jsonErr
	}
	return ret, nil
}

// GetUser will retrieve details about a user from GitHub
func (c *Client) GetUser(login string) (User, error) {
	var ret User
	uri := viper.GetString("github_api") + "users/" + login
	resp, err := c.doGet(uri, nil)
	if err != nil {
		return ret, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ret, errorFromResponse(resp)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	jsonErr := json.Unmarshal(body, &ret)
	if jsonErr != nil {
		return ret, jsonErr
	}
	return ret, nil
}

// -----------------------------------------------
// Most stuff below here is pretty boilerplate

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient() (*Client, error) {
	var c http.Client
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	c.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	return &Client{httpClient: &c}, nil
}
func (c *Client) doDelete(uri string, body interface{}) (*http.Response, error) {
	return c.doMethod("DELETE", uri, body)
}
func (c *Client) doGet(uri string, body interface{}) (*http.Response, error) {
	return c.doMethod("GET", uri, body)
}
func (c *Client) doPatch(uri string, body interface{}) (*http.Response, error) {
	return c.doMethod("POST", uri, body)
}
func (c *Client) doPost(uri string, body interface{}) (*http.Response, error) {
	return c.doMethod("POST", uri, body)
}
func (c *Client) doPut(uri string, body interface{}) (*http.Response, error) {
	return c.doMethod("POST", uri, body)
}
func (c *Client) doMethod(method string, uri string, body interface{}) (*http.Response, error) {
	var err error
	var jsonBody []byte
	var req *http.Request
	if body != nil {
		jsonBody, err = json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("could not marshall json body: %v", err)
		}
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}
	if err != nil {
		return nil, err
	}

	// Set Github Authorization Token in the Header
	tokenvalue := "token " + viper.GetString("github_token")
	req.Header.Set("Authorization", tokenvalue)

	return c.httpClient.Do(req)
}

func errorFromResponse(resp *http.Response) error {

	statusCode := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%d %s", statusCode, string(body))
}

func writePartField(w *multipart.Writer, fieldname, value, contentType string) error {
	p, err := createFormField(w, fieldname, contentType)
	if err != nil {
		return err
	}
	_, err = p.Write([]byte(value))
	return err
}

// quoteEscaper replaces some special characters in a given string.
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// escapeQuotes replaces single quotes and double-backslashes in the current string.
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// createFormField creates the MIME field for a POST request.
func createFormField(w *multipart.Writer, fieldname, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
