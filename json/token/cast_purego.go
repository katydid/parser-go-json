//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

//go:build purego

package token

import (
	"encoding/binary"
	"math"
)

func castToInt64(bs []byte) int64 {
	return binary.LittleEndian.Uint64(bs)
}

func castFromInt64(i int64, alloc func(size int) []byte) []byte {
	bs := alloc(8)
	binary.LittleEndian.PutUint64(bs, uint64(i))
	return bs
}

func castToFloat64(bs []byte) float64 {
	u := binary.LittleEndian.Uint64(bs)
	return math.Float64frombits(u)
}

func castFromFloat64(f float64, alloc func(size int) []byte) []byte {
	bs := alloc(8)
	u := math.Float64bits(f)
	binary.LittleEndian.PutUint64(bs, u)
	return bs
}
