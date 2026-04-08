// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tag

import (
	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/parse"
)

func translateHint(h jsonparse.Hint) parse.Hint {
	switch h {
	case jsonparse.UnknownHint:
		return parse.UnknownHint
	case jsonparse.ObjectOpenHint, jsonparse.ArrayOpenHint:
		return parse.EnterHint
	case jsonparse.ObjectCloseHint, jsonparse.ArrayCloseHint:
		return parse.LeaveHint
	case jsonparse.KeyHint:
		return parse.FieldHint
	case jsonparse.ValueHint:
		return parse.ValueHint
	}
	panic("unreachable")
}
