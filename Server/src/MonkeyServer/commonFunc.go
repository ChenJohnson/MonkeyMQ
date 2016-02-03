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
	"errors"
	//	"log"
)

func getIntFromBytes(bytes []byte) (int, error) {
	if len(bytes) != 4 {
		return -1, errors.New("wrong length")
	}

	iResult := int(bytes[0])*256*256*256 + int(bytes[1])*256*256 + int(bytes[2])*256 + int(bytes[3])
	return iResult, nil
}

func getBytesFromInt32(iLen int) []byte {
	retBytes := make([]byte, 4)

	retBytes[0] = byte(iLen / (256 * 256 * 256))
	retBytes[1] = byte(iLen % (256 * 256 * 256) / (256 * 256))
	retBytes[2] = byte(iLen % (256 * 256) / 256)
	retBytes[3] = byte(iLen % 256)

	return retBytes
}

func combineBytes(first []byte, iFirst int, second []byte, iSecond int) []byte {
	retBytes := make([]byte, iFirst+iSecond)
	for i := 0; i < iFirst; i++ {
		retBytes[i] = first[i]
	}

	for j := 0; j < iSecond; j++ {
		retBytes[iFirst+j] = second[j]
	}

	return retBytes
}
