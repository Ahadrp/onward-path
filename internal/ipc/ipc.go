package ipc

import (
	"bytes"
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	netURL "net/url"
	"strings"

	"net/http/cookiejar"
)

var (
	HTTP_PORT = "2332"
)

type IPC struct {
}

func New() *IPC {
	return &IPC{}
}

func (i IPC) Load() error {
	log.Println("IPC module has been loaded")
	return nil
}

func (i IPC) Run() error {
	// go i.tcpListen()

	log.Printf("Listening to http (2332 port)")
	httpErrCh := make(chan error)
	go i.httpListen(httpErrCh)
	go func() {
		if err := <-httpErrCh; err != nil {
			log.Panicf("Error in listening to http (2332 port): ", err)
		}
	}()

	log.Println("IPC module has been run")
	return nil
}

func (i IPC) tcpListen() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Printf("Couldn't listen to tcp port: ", err)
		return
	}
	log.Printf("Listening on 2000...")
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error in accepting connection: ", err)
			return
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

func (i IPC) httpListen(errCh chan error) {
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		errCh <- err
	}
}

func PostLogin(_url string, formData map[string]string, cookie *cookiejar.Jar) (string, error) {
	url := "http://" + _url

    if cookie == nil {
		return "", fmt.Errorf("PostLogin | cookie can not be null")
    }

	// Prepare form values
	data := netURL.Values{}
	for k, v := range formData {
		data.Set(k, v)
	}

	// Create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Create POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating request | URL: '%s' | Error: '%s'", url, err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request | URL: '%s' | Error: '%s'", url, err)
		return "", err
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response | URL: '%s' | Error: '%s'", url, err)
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Unexpected status code | URL: '%s' | Code: %d", url, resp.StatusCode)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// set cookie
	*cookie = *jar

	log.Printf("PostLogin success | URL: '%s' | Response: '%s'", url, string(body))

	return string(body), nil
}

func Post(_url string, data string, cookie *cookiejar.Jar) (string, error) {
	url := "http://" + _url

    if cookie == nil {
		return "", fmt.Errorf("Post | cookie can not be null")
    }

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar: cookie,
	}

	// Prepare request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Printf("Error creating request | URL: '%s' | Data: '%s' | Error: '%s'", url, data, err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request using the custom client
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending POST request | URL: '%s' | Data: '%s' | Error: '%s'", url, data, err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading POST response | URL: '%s' | Data: '%s' | Status code: '%d' | Error: '%s'", url, data, resp.StatusCode, err)
		return "", err
	}

	return string(body), nil
}
