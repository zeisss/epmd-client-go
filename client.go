package epmd

import (
	"bytes"
	"encoding/binary"
	"net"

//  "log"
)

var NAMES_REQ byte = byte('n')
var PORT_PLEASE2_REQ byte = byte('z')
var PORT2_RESP = byte('w')

const (
	S_DEFAULT_PORT = "4369"
	NL             = byte(10)
	SP             = byte(32)
)

type Name struct {
	Name string
	Port string
}

type NodeInfo struct {
	NodeName       string
	Port           uint16
	NodeType       rune // H = Hidden or M
	Protocol       byte
	HighestVersion uint16
	LowestVersion  uint16
	Extra          []byte
}

type Client struct {
	hostname string
}

func (this *Client) Names() ([]Name, error) {
	return Names(this.hostname)
}

func (this *Client) Get(name string) (*NodeInfo, error) {
	return Get(this.hostname, name)
}

// Returns a client for the local running epmd
func Local() (*Client, error) {
	return NewClient("localhost")
}

// Creates a epmd client connection.
// 
func NewClient(hostname string) (*Client, error) {
	return &Client{hostname}, nil
}

func Names(epmd_hostname string) ([]Name, error) {
	// Connect
	laddr := net.JoinHostPort(epmd_hostname, S_DEFAULT_PORT)
	conn, err := net.Dial("tcp", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send Request
	n, err := sendRequest(conn, NAMES_REQ, nil)
	if err != nil {
		return nil, err
	}

	// Read Response
	var buffer []byte = make([]byte, 1024*1024)
	n, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Parse
	return parseNamesResponse(buffer[0:n])
}

func Get(epmd_hostname, nodeName string) (*NodeInfo, error) {
	// Connect
	laddr := net.JoinHostPort(epmd_hostname, S_DEFAULT_PORT)
	conn, err := net.Dial("tcp", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Request
	n, err := sendRequest(conn, PORT_PLEASE2_REQ, []byte(nodeName))
	if err != nil {
		return nil, err
	}

	// Read Response
	var buffer []byte = make([]byte, 1024*1024)
	n, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Parsing
	return parsePort2PleaseResponse(buffer[0:n])
}

func sendRequest(conn net.Conn, reqId byte, data []byte) (int, error) {
	var req []byte = make([]byte, 3)
	binary.BigEndian.PutUint16(req[0:2], uint16(1+len(data)))
	req[2] = reqId
	req = bytes.Join([][]byte{req, data}, []byte{})
	// log.Printf("REQ: %v\n", req)
	return conn.Write(req)
}
