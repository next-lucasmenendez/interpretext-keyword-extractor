[![GoDoc](https://godoc.org/github.com/next-lucasmenendez/interpretext-keyword-extractor?status.svg)](https://godoc.org/github.com/next-lucasmenendez/interpretext-keyword-extractor)
[![Report](https://goreportcard.com/badge/github.com/next-lucasmenendez/gotagger)](https://goreportcard.com/report/github.com/next-lucasmenendez/interpretext-keyword-extractor)

# Interpretext Keyword Extractor
Simple keyword extraction

## Installation
```bash
go install github.com/next-lucasmenende/interpretext-keyword-extraction
```

### Stopwords
If you want to extend stopword list, create a file, named as language code, into a any folder (for example: `en` file will contain English stopwords). Then, set env var `STOPWORDS` to that folder path.
Extended stopword lists can be found in [Stopwords ISO profile](https://github.com/stopwords-iso).

## Demo
```go
package main

import (
    "fmt"
    tokenizer "github.com/next-lucasmenendez/interpretext-tokenizer"
    keywords "github.com/next-lucasmenendez/interpretext-keyword-extractor"
)

func main() {
    var limit int = 15
    var lang string = "<input-lang>"
    var text string = "<input-text>"
    
    var words [][]string
    for _, s := range tokenizer.Sentences(text) {
        words = append(words, tokenizer.Words(s))
    }
    
    if tags, err := keywords.GetTags(words, lang, limit); err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("%q\n", tags)
    }
}
```
