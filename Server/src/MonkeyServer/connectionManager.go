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
	"log"
	"net"
	"regexp"
)

type ConnUnit struct {
	IP       string
	PubConn  *net.Conn
	SubConn  *net.Conn
	SubTopic string
	SubChann chan Message
}

func NewConnUnit(ip string) *ConnUnit {
	connUnit := new(ConnUnit)
	connUnit.IP = ip
	connUnit.SubTopic = ""

	return connUnit
}

func (pConnUnit *ConnUnit) SetPubConn(pubConn *net.Conn) {
	pConnUnit.PubConn = pubConn
}

func (pConnUnit *ConnUnit) SetSubConn(subConn *net.Conn, strSub string) {
	pConnUnit.SubConn = subConn

	pConnUnit.SubTopic = strSub
}

type ConnManager struct {
	allConns       []*ConnUnit
	topic2ConnUnit map[string][]*ConnUnit
	ip2ConnUnit    map[string]*ConnUnit

	bCanSet chan bool
}

func NewConnManager() *ConnManager {
	ret := new(ConnManager)
	ret.allConns = make([]*ConnUnit, 512)
	ret.topic2ConnUnit = make(map[string][]*ConnUnit, 512)
	ret.ip2ConnUnit = make(map[string]*ConnUnit, 512)

	ret.bCanSet = make(chan bool, 1)

	return ret
}

func (pConnManager *ConnManager) AddConn(connObj *ConnUnit) {
	pConnManager.bCanSet <- true

	if _, ok := pConnManager.ip2ConnUnit[connObj.IP]; !ok {
		pConnManager.ip2ConnUnit[connObj.IP] = connObj
		pConnManager.allConns = append(pConnManager.allConns, connObj)
	}

	_ = <-pConnManager.bCanSet
}

func (pConnManager *ConnManager) AddPubConn(ip string, pubConn *net.Conn) {
	pConnManager.bCanSet <- true

	if _, ok := pConnManager.ip2ConnUnit[ip]; !ok {
		newConnUnit := NewConnUnit(ip)
		pConnManager.allConns = append(pConnManager.allConns, newConnUnit)
		pConnManager.ip2ConnUnit[ip] = newConnUnit
	}
	pConnManager.ip2ConnUnit[ip].PubConn = pubConn

	<-pConnManager.bCanSet
}

func (pConnManager *ConnManager) AddSubConn(ip string, subConn *net.Conn, strSub string) {
	pConnManager.bCanSet <- true

	if _, ok := pConnManager.ip2ConnUnit[ip]; !ok {
		newConnUnit := NewConnUnit(ip)
		pConnManager.allConns = append(pConnManager.allConns, newConnUnit)
		pConnManager.ip2ConnUnit[ip] = newConnUnit
	}

	pConnManager.ip2ConnUnit[ip].SetSubConn(subConn, strSub)

	<-pConnManager.bCanSet
}

func (pConnManager *ConnManager) AddPubTopic(strPublish string) {
	pConnManager.bCanSet <- true

	for _, v := range pConnManager.allConns {
		if v != nil && v.SubTopic != "" {
			if ok, _ := regexp.MatchString(v.SubTopic, strPublish); ok {
				pConnManager.topic2ConnUnit[strPublish] = append(pConnManager.topic2ConnUnit[strPublish], v)
			}
		}

	}

	<-pConnManager.bCanSet
}

// write message to related channels subscribed the topic
func (pConnManager *ConnManager) Write2Channels(strTopic string, msg Message) {
	pConnManager.bCanSet <- true

	lst := pConnManager.topic2ConnUnit[strTopic]

	<-pConnManager.bCanSet

	//log.Println("lst size:", len(lst))
	for _, v := range lst {
		v.SubChann <- msg
	}
}
