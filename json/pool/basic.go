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

package pool

type pool struct {
	free [][]byte
	busy [][]byte
}

func New() Pool {
	return &pool{
		free: make([][]byte, 0),
		busy: make([][]byte, 0),
	}
}

func (p *pool) FreeAll() {
	p.free = append(p.free, p.busy...)
	p.busy = p.busy[:0]
}

func (p *pool) Alloc(size int) []byte {
	for i := 0; i < len(p.free); i++ {
		if len(p.free[i]) >= size {
			buf := p.free[i]
			p.free[i] = p.free[len(p.free)-1]
			p.free = p.free[:len(p.free)-1]
			p.busy = append(p.busy, buf)
			return buf[:size]
		}
	}
	// always allocate a big buffer, so hits when searching are very likely
	buf := make([]byte, min(size*2, 1000))
	p.busy = append(p.busy, buf)
	return buf[:size]
}
