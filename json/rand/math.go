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

package rand

import (
	mathrand "math/rand"
	"time"
)

type Rand interface {
	// Intn returns a pseudo-random number in [0,n).
	Intn(int) int
	// Seed returns the seed used to create the random number generator.
	Seed() int64
}

type rand struct {
	*mathrand.Rand
	seed int64
}

// Seed returns the seed used to create the random number generator.
func (r *rand) Seed() int64 {
	return r.seed
}

// NewRand returns a new Rand with random seed.
func NewRand() Rand {
	seed := time.Now().UnixNano()
	return NewRandWithSeed(seed)
}

// NewRandWithSeed allows the seed to be manually set for Rand.
func NewRandWithSeed(seed int64) Rand {
	return &rand{
		Rand: mathrand.New(mathrand.NewSource(seed)),
		seed: seed,
	}
}
