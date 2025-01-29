## Fast Float

fastfloat has been originally copied from https://github.com/valyala/fastjson
This included files:
* parse.go
* parse_test.go

Changes include:
* BestEffort functions and tests have been removed, since we do not need them.
* Inline errors were replaced with predefined errors to avoid allocations.
* Functions now take `[]byte` as input, instead of `string`.
* Checking for inifinity and NaN is done with manually written functions instead of using the strings package.
