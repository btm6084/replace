Simple find and replace. Searches every file in a given directory, and all sub directories. Ignores .git directory entirely.

# Install
```go get github.com/btm6084/replace```

# Usage

replace "search term" "replace term" path/to/directory extension_filter

Path defaults to ./

Extension filter defaults to no filter

Example
```
replace "fmt.Prinln" "fmt.Println" ~/golang/src/github.com/btm6084/replace .go
```
