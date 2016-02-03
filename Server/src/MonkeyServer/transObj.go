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

type Message struct {
	Len     int    `json:"len"`
	Content []byte `json:"content"`
}

// publish object class
type PubMessage struct {
	Topic string  `json:"topic"`
	Msg   Message `json:"message"`
}

// tans object when creating connection
type ConnectionApply struct {
	Kind  string `json:"kind"`
	Op    string `json:"op"`
	Topic string `json:"topic"`
}

// return object when creating connection
type ConnectionResponse struct {
	Kind  string `json:"kind"`
	Value bool   `json:"value"`
}

// trans object
type TransObject struct {
	Kind    string  `json:"kind"`
	Op      string  `json:"op"`
	Topic   string  `json:"topic"`
	Content Message `json:"content"`
}
