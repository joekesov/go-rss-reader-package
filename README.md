# github.com/joekesov/go-rss-reader-package

## Usage

### Initialize your module

```bash
go mod init example.com/my-golib-demo
```

### Get the go-rss-reader-package

```bash
go get github.com/joekesov/go-rss-reader-package@v0.1.0
```

And in your main function is it like below

```go
package main

import (
    ...
    "github.com/joekesov/go-rss-reader-package"
)

func main() {
	urls := []string {"https://example.com/rss/feed"}
	rssItems, err := jreader.Parse(urls)
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
