## parser-go-json

Parser for JSON in Go

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