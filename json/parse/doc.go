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
//	p.Next() // k
//	p.Token() // StringKind, "num"
//
//	p.Next() // v
//	p.Token() // Float64Kind, 3.14
//
//	p.Next() // k
//	p.Token() // StringKind, "arr"
//
//	p.Next() // [
//
//	p.Next() // v
//	p.Token() // NullKind
//
//	p.Next() // v
//	p.Token() // FalseKind
//
//	p.Next() // v
//	p.Token() // TrueKind
//
//	p.Skip()
//
//	p.Next() // k
//	p.Token() // StringKind, "obj"
//
//	p.Next() // {
//
//	p.Next() // k
//	p.Token() // StringKind, "k"
//
//	p.Next() // v
//	p.Token() // StringKind, "v"
//
//	p.Next() // k
//	p.Token() // StringKind, "a"
//
//	p.Skip()
//
//	p.Skip()
//
//	p.Next() // }
package parse
