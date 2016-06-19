package ss13

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

// bytes 1 to 2 = unknown, magic number?
// byte 3 = unknown, part of byte 4?
// byte 4 = length of request (5 null bytes + len(?players) + last null byte = 14 = 0x0e)
// bytes 5 to 9 = unknown, padding?
// bytes 10 to 17 = request (?players)
// byte 18 = end of line suffix (null byte)
var QUERY_PLAYERS = []byte{0x00, 0x83, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x00}

func poll_players(host string, timeout int) (int, error) {
	// Connect to the server
	dur := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", host, dur)
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(dur))

	// Send the ?players query
	size, err := conn.Write(QUERY_PLAYERS)
	if err != nil {
		return -1, err
	}
	if size != len(QUERY_PLAYERS) {
		return -1, fmt.Errorf("Failed to send all data")
	}

	// Read the server's response (should be 9 bytes long)
	buf := make([]byte, 9)
	size, err = conn.Read(buf)
	if err != nil {
		return -1, err
	}
	if size != 9 || !bytes.HasPrefix(buf, []byte{0x00, 0x83, 0x00, 0x05, 0x2a}) {
		return -1, fmt.Errorf("Received invalid response")
	}

	// Grab the encoded value, which should be a float32 (4 bytes long)
	// buf[5:] = last 4 bytes
	var players float32
	err = binary.Read(bytes.NewReader(buf[5:]), binary.LittleEndian, &players)
	if err != nil {
		return -1, err
	}

	return int(players), nil
}

func (i *Instance) PollServers(servers []ServerConfig, timeout int) []*RawServerData {
	var wg sync.WaitGroup
	var tmp []*RawServerData
	for _, s := range servers {
		wg.Add(1)
		go func(s ServerConfig) {
			defer wg.Done()
			players, e := poll_players(s.GameUrl, timeout)
			if e != nil {
				Log("Error polling server %s: %s", s.GameUrl, e)
				return
			}
			gameurl := fmt.Sprintf("byond://%s", s.GameUrl)
			// TODO: data race
			tmp = append(tmp, &RawServerData{s.Title, gameurl, s.SiteUrl, players, Now()})
		}(s)
	}
	wg.Wait()
	return tmp
}
