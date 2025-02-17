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

type config struct {
	// maxDepth is the maximum number of objects/arrays to generate down. It makes sure the algorithm terminates.
	maxDepth        int
	maxObjectFields int
	maxArrayLength  int
	maxStringLength int
	numberEdgeCases []string
	// r.Intn(c.numberEdgeCaseOdds) == 0 will result in a random edge case being generated.
	numberEdgeCaseOdds int
	maxSpaces          int
}

type Option func(*config)

func newConfig(opts ...Option) *config {
	// default Config
	c := &config{
		maxDepth:           5,
		maxObjectFields:    10,
		maxArrayLength:     10,
		maxStringLength:    100,
		numberEdgeCases:    defaultNumberEdgeCases,
		numberEdgeCaseOdds: 100,
		maxSpaces:          5,
	}
	// apply options
	for _, o := range opts {
		o(c)
	}
	return c
}

var defaultNumberEdgeCases = []string{
	"9223372036854775807",                                // math.MaxInt64
	"9223372036854775808",                                // math.MaxInt64 + 1
	"-9223372036854775808",                               // math.MinInt64
	"-9223372036854775809",                               // math.MinInt64 - 1
	"18446744073709551615",                               // math.MaxUint64
	"18446744073709551616",                               // math.MaxUint64 + 1
	"1.79769313486231570814527423731704356798070e+308",   // math.MaxFloat64
	"2.79769313486231570814527423731704356798070e+308",   // > math.MaxFloat64
	"4.9406564584124654417656879286822137236505980e-324", // math.SmallestNonzeroFloat64
}

func WithMaxDepth(maxDepth int) Option {
	return func(c *config) {
		c.maxDepth = maxDepth
	}
}

func WithMaxObjectFields(maxObjectFields int) Option {
	return func(c *config) {
		c.maxObjectFields = maxObjectFields
	}
}

func WithMaxArrayLength(maxArrayLength int) Option {
	return func(c *config) {
		c.maxArrayLength = maxArrayLength
	}
}

func WithMaxStringLength(maxStringLength int) Option {
	return func(c *config) {
		c.maxStringLength = maxStringLength
	}
}

// WithMaxSpaces sets the maximum number of spaces generated at each opportunity to generate spaces.
func WithMaxSpaces(maxSpaces int) Option {
	return func(c *config) {
		c.maxSpaces = maxSpaces
	}
}
