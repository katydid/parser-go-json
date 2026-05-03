## parser-go-json

Parser for JSON in Go, which tries to not allocate any memory.

We can parse json without unmarshaling it into a Go `struct` using the [Parser interface](https://github.com/katydid/parser-go/parse):

```go
import "github.com/katydid/parser-go-json/json"

func main() {
    jsonString := `{"otherfield": 123, "myfield": "myvalue"}`
	jsonParser := json.NewParser()
	jsonParser.Init([]byte(jsonString))
	myvalue, err := GetMyField(jsonParser)
	if err != nil {
		panic(err)
	}
	println(myvalue)
}
```

We can then use the parser to decode only `myfield` and skip over other fields and return `"myvalue"`:

```go
import (
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)


func GetMyField(p parse.Parser) (string, error) {
	hint, err := p.Next()
	if hint != parse.EnterHint {
		return "", errors.New("expected object")
	}
	if err != nil {
		return "", err
	}
	for {
		hint, err = p.Next()
		if err != nil {
			return "", err
		}
		if hint != parse.FieldHint {
			return "", errors.New("expected field")
		}
		kind, fieldName, err := p.Token()
		if err != nil {
			return "", err
		}
		if kind != parse.StringKind {
			return "", errors.New("expected string")
		}
		if cast.ToString(fieldName) == "myfield" {
			hint, err = p.Next()
			if err != nil {
				return "", err
			}
			if hint != parse.ValueHint {
				return "", errors.New("expected field")
			}
			kind, val, err := p.Token()
			if err != nil {
				return "", err
			}
			if kind != parse.StringKind {
				return "", errors.New("expected string")
			}
			return cast.ToString(val), nil
		} else {
			p.Skip()
		}
	}
}
```

## Special Considerations

* The parser uses a buffer pool, which will allocate memory until it is warmed up.
* Buffers are reused and pooled. This means that the `String` and `Bytes` methods, returns a `string` and `[]byte` respectively that should be copied if it is needed again before calling `Next`, `Token` or `Skip`.
* Arrays are indexed, which means that `["a","b","c"]` will be parsed into something that looks like an integer indexed map: `{O: "a", 1: "b", 2: "c"}`.

## Thank you

Thanks to the following people for consulting on the project:

- [Jacques Marais](https://www.linkedin.com/in/ajacquesmarais/)
- [Johan Brandhorst-Satzkorn](https://www.linkedin.com/in/jbrandhorst/)

I still made all the bad design decisions, so don't blame them.
They only made the project better than it would have been without their advice.
