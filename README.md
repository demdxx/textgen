Textgen
=======

Library text generation.

[![Gratipay](http://img.shields.io/gratipay/demdxx.svg)](https://gratipay.com/demdxx/)

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
