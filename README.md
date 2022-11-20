# github.com/mitchallen/go-lib

## Usage

### Initialize your module

```
$ go mod init example.com/my-golib-demo
```

### Get the go-lib module

Note that you need to include the **v** in the version tag.

```
$ go get github.com/joekesov/go-rss-reader-package@v0.1.0
```

```go
package main

import (
    ...
    jrss "github.com/joekesov/go-rss-reader-package"
)

func main() {
	urls := []string {"https://example.com/rss/feed"}
	rssItems, err := jrss.Parse(urls)
	if err != nil {
		// handle error
	}
	...
}
```

## Testing

```
$ go test
```

## Tagging

```
$ git tag v0.1.0
$ git push origin --tags
```
