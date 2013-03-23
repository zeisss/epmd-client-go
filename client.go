package epmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

//	"log"
)

var NAMES_REQ byte = byte('n')
var PORT_PLEASE2_REQ byte = byte('z')
var PORT2_RESP = byte('w')

const (
	ALIVE2_REQ  byte = 120
	ALIVE2_RESP byte = 121

	S_DEFAULT_PORT = "4369"
	NL             = byte(10)
	SP             = byte(32)

	NODE_TYPE_HIDDEN byte = byte(72)
	NODE_TYPE_ERLANG byte = byte(77)
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

// Register the current process to the locally running epmd.
// The protocol defaults to 0=tcp/ip
//
// @return an error if the registration failed (either because no EPMD is running or the registration failed)
func Register(port uint16, node_type byte, highest_version, lowest_version uint16, name, extra string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ":"+S_DEFAULT_PORT)
	if err != nil {
		return nil, err
	}

	// Build ALIVE2_REQ
	var req []byte = make([]byte, 0, 12+len(name)+len(extra))
	var buf *bytes.Buffer = bytes.NewBuffer(req)
	binary.Write(buf, binary.BigEndian, port)
	buf.WriteByte(node_type)
	buf.WriteByte(0)
	binary.Write(buf, binary.BigEndian, highest_version)
	binary.Write(buf, binary.BigEndian, lowest_version)
	binary.Write(buf, binary.BigEndian, uint16(len(name)))
	buf.WriteString(name)
	binary.Write(buf, binary.BigEndian, uint16(len(extra)))
	buf.WriteString(extra)
	n, err := sendRequest(conn, ALIVE2_REQ, buf.Bytes())

	// Read response
	var resp []byte = make([]byte, 4)
	n, err = conn.Read(resp)
	if err != nil {
		return nil, err
	}
	if n != 4 || resp[0] != ALIVE2_RESP {
		return nil, fmt.Errorf("EPMD returned unexpected code %v", resp[0])
	}
	if resp[1] > 0 {
		return nil, fmt.Errorf("EPMD: registration failed with error code %v", int(resp[1]))
	}
	return conn, nil
}

func sendRequest(conn net.Conn, reqId byte, data []byte) (int, error) {
	var req []byte = make([]byte, 3)
	binary.BigEndian.PutUint16(req[0:2], uint16(1+len(data)))
	req[2] = reqId
	req = bytes.Join([][]byte{req, data}, []byte{})
	//	log.Printf("REQ: %v\n", req)
	return conn.Write(req)
}
