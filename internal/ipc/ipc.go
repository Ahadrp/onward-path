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

func PostLogin(_url string, formData map[string]string) (string, error) {
	url := "http://" + _url

	// Prepare form values
	data := netURL.Values{}
	for k, v := range formData {
		data.Set(k, v)
	}

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error in sending post request | URL: '%s' | Data: '%s' | Error: '%s'", url, data, err)
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in reading post request | URL: '%s' | Data: '%s' | Error: '%s'", url, data, err)
		return "", err
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Error in reading post request | URL: '%s' | Data: '%s' | Status code: '%d' | Error: '%s'", url, data, resp.StatusCode, err)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("Sending post request was successful! | URL: '%s' | Data: '%s' | Response: '%s'", url, data, string(body))
	return string(body), nil
}

func Post(_url string, data string) (string, error) {
	url := "http://" + _url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Printf("Error in sending post request | URL: '%s' | Data: '%s' | Status code: '%d' | Error: '%s'", url, data, resp.StatusCode, err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in reading post request | URL: '%s' | Data: '%s' | Status code: '%d' | Error: '%s'", url, data, resp.StatusCode, err)
		return "", err
	}

	return string(body), nil
}

