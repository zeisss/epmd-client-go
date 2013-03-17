package epmd

import (
	"bytes"
	"encoding/binary"
	//  "log"
	"fmt"
)

// [119       # Response code (ascii)
//    0       # Result code (0 = ok, >0 = error)
//  221 86    # Port (uint16, big endian)
//   77       # NodeType (M or H = hidden node)
//    0       # Protocol
//  0 5       # Highest version 
//  0 5       # Lowest version
//  0 6       # Length of Name field
//  114 97 98 98 105 116 # Name field
//   0        # Length of extra field
//   0]       # ?????????????
func parsePort2PleaseResponse(data []byte) (*NodeInfo, error) {
	if data[0] != PORT2_RESP {
		return nil, fmt.Errorf("Expected response code PORT2_RESP")
	}
	if data[1] != 0 {
		return nil, fmt.Errorf("Epmd responded with error code %v", int(data[1]))
	}

	// log.Printf("RES: %v\n", data)
	var nodeInfo NodeInfo
	nodeInfo.Port = binary.BigEndian.Uint16(data[2:4])
	nodeInfo.NodeType = rune(data[4])
	nodeInfo.Protocol = data[5]
	nodeInfo.HighestVersion = binary.BigEndian.Uint16(data[6:8])
	nodeInfo.LowestVersion = binary.BigEndian.Uint16(data[8:10])

	nameLen := binary.BigEndian.Uint16(data[10:12])
	nodeInfo.NodeName = string(data[12 : 12+nameLen])

	data = data[12+nameLen:]
	extraLen := binary.BigEndian.Uint16(data[0:2])
	nodeInfo.Extra = data[2 : 2+extraLen]

	return &nodeInfo, nil
}

// [0 0 17 17  # EPMD Port
//
//  110 97 109 101 32 102 111 111 32 97 116 32 112 111 114 116 32 53 49 51 54 56 10 
//  n   a  m   e   SP f   o   o   SP A  T   SO p   o   r   t   SP 5  1  3  6  8  NL
//
//  110 97 109 101 32 114 97 98 98 105 116 32 97 116 32 112 111 114 116 32 53 54 54 54 50 10
//  110 97 109 101 32 114 105 97 107 32 97 116 32 112 111 114 116 32 52 57 50 50 51 10
// ]
func parseNamesResponse(data []byte) ([]Name, error) {
	data = data[4:] // Skip the first four bytes (which is the epmd port)

	bnodes := bytes.Split(data, []byte{NL})
	nodes := make([]Name, len(bnodes)-1)

	for i, bnode := range bnodes {
		if len(bnode) == 0 {
			break
		}
		bfields := bytes.Split(bnode, []byte{SP})

		nodes[i].Name = string(bfields[1])
		nodes[i].Port = string(bfields[4])
	}
	return nodes, nil
}
