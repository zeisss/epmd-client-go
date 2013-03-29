package epmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	ALIVE2_REQ  byte = 120
	ALIVE2_RESP byte = 121

	NAMES_REQ        byte = 110
	PORT_PLEASE2_REQ byte = 122
	PORT2_RESP       byte = 119

	DEFAULT_PORT = "4369"
	NL           = byte(10)
	SP           = byte(32)

	NODE_TYPE_HIDDEN byte = byte(72)
	NODE_TYPE_ERLANG byte = byte(77)
)

var epmd_port = DEFAULT_PORT

type Name struct {
	Name string
	Port string
}

type NodeInfo struct {
	NodeName       string
	Port           uint16
	NodeType       byte
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

func SetPort(new_port string) {
	epmd_port = new_port
}

func Names(epmd_hostname string) ([]Name, error) {
	// Connect
	laddr := net.JoinHostPort(epmd_hostname, epmd_port)
	conn, err := net.Dial("tcp", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send Request
	_, err = sendRequest(conn, NAMES_REQ, nil)
	if err != nil {
		return nil, err
	}

	// Read Response

	var data []byte = make([]byte, 0, 16*1024)
	buffer := bytes.NewBuffer(data)
	n2, err := io.Copy(buffer, conn)
	if err != nil {
		return nil, err
	}

	// Parse
	return parseNamesResponse(buffer.Bytes()[0:n2])
}

// Request detailed data about the node with the given name
// from the (remote) epmd_hostname.
func Get(epmd_hostname, nodeName string) (*NodeInfo, error) {
	// Connect
	laddr := net.JoinHostPort(epmd_hostname, epmd_port)
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
// Returns an error if the registration failed (either because no EPMD is running or the registration failed)
func Register(port uint16, node_type byte, highest_version, lowest_version uint16, name, extra string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ":"+epmd_port)
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
