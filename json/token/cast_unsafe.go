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

//go:build !purego

package token

import (
	"math"
	"reflect"
	"unsafe"
)

func castToInt64(bs []byte) int64 {
	return *(*int64)(unsafe.Pointer(&bs[0]))
}

func castFromInt64(i int64, _alloc func(size int) []byte) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&i)), 8)
}

func castToFloat64(bs []byte) float64 {
	u := *(*uint64)(unsafe.Pointer(&bs[0]))
	return math.Float64frombits(u)
}

func castFromFloat64(f float64, _alloc func(size int) []byte) []byte {
	u := math.Float64bits(f)
	return unsafe.Slice((*byte)(unsafe.Pointer(&u)), 8)
}

func castToString(buf []byte) string {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	strHeader := reflect.StringHeader{Data: header.Data, Len: header.Len}
	return *(*string)(unsafe.Pointer(&strHeader))
}
