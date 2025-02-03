## parser-go-json

Parser for JSON in Go, which tries to not allocate any memory.

We can parse json without unmarshaling it into a Go `struct` using the [Parser interface](https://github.com/katydid/parser-go):

```go
import "github.com/katydid/parser-go-json/json"

func main() {
    jsonString := `{"otherfield": 123, "myfield": "myvalue"}`
    jsonParser := json.NewJsonParser()
    if err := jsonParser.Init([]byte(jsonString)); err != nil {
        panic(err)
    }
    myvalue, err := GetMyField(jsonParser)
    if err != nil {
        panic(err)
    }
    println(value)
}
```

We can then use the parser to decode only `myfield` and skip over other fields and return `"myvalue"`:

```go
func GetMyField(p parser.Interface) (string, error) {
	for {
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		fieldName, err := p.String()
		if err != nil {
			return "", err
		}
		if fieldName != "myfield" {
			continue
		}
		p.Down()
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		return p.String()
	}
	return "", nil
}
```

## Special Considerations

* The parser uses a buffer pool, which will allocate memory until it is warmed up.
* Buffers are reused and pooled. This means that the `String` and `Bytes` methods, returns a `string` and `[]byte` respectively that should be copied if it is needed again after calling `Next`, `Up` or `Down`.
* Arrays are indexed, which means that `["a","b","c"]` will be parsed something that looks like an integer indexed map: `{O: "a", 1: "b", 2: "c"}`.