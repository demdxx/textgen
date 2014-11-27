Textgen
=======

Library text generation.

<a href="https://gratipay.com/demdxx/" target="_blank"><img src="//img.shields.io/gratipay/demdxx.svg"></a>

    @copyright Dmitry Ponomarev 2014 <demdxx@gmail.com>
    @license MIT

Example
=======

```go

text := "[Test|Text|Guest|Start it|Crash it] [for|not for] [example|fun|you|us]! Oh ya!"

for s := range textgen.Generate(text, 0) {
  fmt.Println(s+"\n")
}

```
