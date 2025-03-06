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

package token

import (
	"reflect"
	"unsafe"
)

func castToInt64(bs []byte) int64 {
	return *(*int64)(unsafe.Pointer(&bs[0]))
}

func castFromInt64(i int64) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  8,
		Cap:  8,
		Data: uintptr(unsafe.Pointer(&i)),
	}))
}

func unsafeCastFromInt64(i int64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&i)), 8)
}
