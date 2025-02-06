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

//	{"num":3.14,"arr":[null,false,true,1,2],"obj":{"k":"v","a":[1,2,3],"b":1,"c":2}}
//
// Can be parsed using Next, Skip and tokenizer methods.
//
//	p.Next() // {
//
//	p.Next() // "
//	p.String() // "num"
//
//	p.Next() // 0
//	p.Int() // token.ErrNotInt
//	p.Uint() // token.ErrNotInt
//	p.Double() // 3.14
//
//	p.Next() // "
//	p.String() // "arr"
//
//	p.Next() // [
//
//	p.Next() // n
//
//	p.Next() // ?
//	p.Bool() // false
//
//	p.Next() // ?
//	p.Bool() // true
//
//	p.Skip()
//
//	p.Next() // "
//	p.String() // "obj"
//
//	p.Next() // {
//
//	p.Next() // "
//	p.String() // "k"
//
//	p.Next() // "
//	p.String() // "v"
//
//	p.Next() // "
//	p.String() // "a"
//
//	p.Skip()
//
//	p.Skip()
//
//	p.Next() // }
package parse
