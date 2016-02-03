// Copyright 2016 The MonkeyMQ Authors. All rights reserved.
//
// Beijing Coding Technology Co. Ltd.
// Author: Johnson Chen
// Email: Johnson.Chen0204@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monkeyserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type ServerManager struct {
	ipAddr  string
	strPort string

	connManger *ConnManager
}

func NewServerManager(port string) *ServerManager {
	ret := new(ServerManager)
	ret.ipAddr = "192.168.1.102"
	ret.strPort = port
	ret.connManger = NewConnManager()

	return ret
}

func (pServer *ServerManager) StartServer() {
	l, err := net.Listen("tcp", pServer.ipAddr+":"+pServer.strPort)
	if err != nil {
		log.Fatal("cannot listen on the specified ip and port!")
		panic("cannot listen on the specified ip and port!")
	}
	defer l.Close()

	fmt.Println("starting server...")
	for {
		conn, err := l.Accept()
		if err != nil {
			panic("wrong in accept function")
		}

		go pServer.handleConnection(conn)
	}
}

// handle connection
func (pServer *ServerManager) handleConnection(conn net.Conn) {
	// get client's ip address
	ipAddr := conn.RemoteAddr()
	// get client ip address value using string
	ipAndPort := strings.Split(ipAddr.String(), ":")

	connObj := NewConnUnit(ipAndPort[0])
	pServer.connManger.AddConn(connObj)

	var iPacketSize int = 0
	var leftBytes []byte = make([]byte, 2048)
	var iLeft int = 0

	for {
		//log.Println("1")
		bytes := make([]byte, 10000)
		//log.Println("2")
		iLen, err := conn.Read(bytes)

		var allBytes []byte
		var iRead int = 0

		if iLeft != 0 {
			allBytes = combineBytes(leftBytes, iLeft, bytes, iLen)
			iLeft = 0
		} else {
			allBytes = bytes[0:iLen]
		}

		if iPacketSize == 0 {
			if len(allBytes) < 4 {
				leftBytes = allBytes[0:len(allBytes)]
				iLeft = len(allBytes)
				continue
			} else if len(allBytes) == 4 {
				iPacketSize, err = getIntFromBytes(allBytes)
				if err != nil {
					log.Fatal("wrong length")
				}
				iLeft = 0
				continue
			} else {
				iPacketSize, err = getIntFromBytes(allBytes[0:4])
				if err != nil {
					log.Fatal("wrong lenght")
				}
				iRead = iRead + 4
			}
		}

		if len(allBytes)-iRead < iPacketSize {
			leftBytes = allBytes[iRead:len(allBytes)]
			iLeft = len(allBytes) - iRead
			break
		} else {
			for {
				content := allBytes[iRead : iRead+iPacketSize]
				iRead = iRead + iPacketSize

				//log.Println(string(content))
				// handle json value
				var requestBody TransObject
				err = json.Unmarshal(content, &requestBody)

				log.Println("requestbogy:", requestBody)

				if err != nil {
					panic("wrong json error")
				}

				pServer.handleReuqest(requestBody, conn, ipAndPort[0])

				if len(allBytes)-iRead < 4 {
					iPacketSize = 0
					leftBytes = allBytes[iRead:len(allBytes)]
					iLeft = len(allBytes) - iRead
					break
				} else if len(allBytes)-iRead >= 4 {
					iPacketSize, err = getIntFromBytes(allBytes[iRead : iRead+4])
					if err != nil {
						log.Fatal("wrong length")
					}
					iRead = iRead + 4

					if iPacketSize > len(allBytes)-iRead {
						leftBytes = allBytes[iRead:len(allBytes)]
						iLeft = len(allBytes) - iRead
						break
					}
				}
			}
		}
	}
}

func (pServer *ServerManager) handleReuqest(requestBody TransObject, conn net.Conn, strIp string) {
	//log.Println(requestBody)
	// create connection before trans object
	if requestBody.Kind == `create_connection` {
		if requestBody.Op == `pub` {
			log.Println("pub create connection")
			if _, ok := pServer.connManger.ip2ConnUnit[strIp]; !ok {
				pServer.connManger.ip2ConnUnit[strIp] = NewConnUnit(strIp)
			}
			pServer.connManger.ip2ConnUnit[strIp].SetPubConn(&conn)
		} else {
			if requestBody.Topic != "" {
				if _, ok := pServer.connManger.ip2ConnUnit[strIp]; !ok {
					pServer.connManger.ip2ConnUnit[strIp] = NewConnUnit(strIp)
				}

				pServer.connManger.ip2ConnUnit[strIp].SetSubConn(&conn, requestBody.Topic)
				pServer.connManger.ip2ConnUnit[strIp].SubChann = make(chan Message, 1000)

				go func() {
					for {
						msg := <-pServer.connManger.ip2ConnUnit[strIp].SubChann
						bytes, err := json.Marshal(msg)
						if err != nil {
							panic("error in marshal msg")
						}

						conn.Write(bytes)
					}
				}()
			} else {
				panic("cannot subscribe null string")
			}
		}
	} else if requestBody.Kind == `message` { // trans object
		if requestBody.Op == `pub` {
			if requestBody.Topic != "" {
				pServer.connManger.AddPubTopic(requestBody.Topic)
				pServer.connManger.Write2Channels(requestBody.Topic, requestBody.Content)
			}
		}
	} else { // "heart-beat package"

	}
}
