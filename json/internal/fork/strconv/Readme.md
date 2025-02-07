# strconv

This was originally copied from the Go standard library's strconv package.
Modifications were made to make it more efficient:
1. Parameter inputs are byte slices instead of strings, to avoid copies.
2. Returned errors are all predefined to avoid allocations.
3. Base is always 10.
4. Bitsize is always 64.

New files:
* atoi_errors.go