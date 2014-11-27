/**
 * textgen SEO library for text generation
 *
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Dmitry Ponomarev <demdxx@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package textgen

import (
  "errors"
  "math/rand"
  "regexp"
  "strings"
)

var (
  defaultProcessor *Processor
  errorNoVariants  = errors.New("No generate variants")
)

type Processor struct {
  collectionRegex *regexp.Regexp
  borderSize      int
  separator       string
}

func MakeProcessor(re string, sep string, borderSize int) *Processor {
  r, _ := regexp.Compile(re)
  return &Processor{collectionRegex: r, separator: sep, borderSize: borderSize}
}

func getDefaultProcessor() *Processor {
  if nil == defaultProcessor {
    defaultProcessor = MakeProcessor("\\[[^\\]]+\\]", "|", 1)
  }
  return defaultProcessor
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Methods
///////////////////////////////////////////////////////////////////////////////

func (p *Processor) PrepareVariants(text string) ([]string, [][]string, error) {
  all := p.collectionRegex.FindAllString(text, -1)
  if nil == all {
    return nil, nil, errorNoVariants
  }

  variants := make([][]string, len(all))
  for i, variant := range all {
    variants[i] = strings.Split(variant[p.borderSize:len(variant)-p.borderSize], p.separator)
  }
  return all, variants, nil
}

func PrepareVariants(text string) ([]string, [][]string, error) {
  return getDefaultProcessor().PrepareVariants(text)
}

func (p *Processor) ProcessRandom(text string, words *[]string) (string, error) {
  all, variants, err := p.PrepareVariants(text)
  if nil != err {
    return "", err
  }
  if nil == variants || len(variants) < 1 {
    return "", errorNoVariants
  }
  result, _, err := ProcessExt(text, all, variants, nil, words)
  return result, err
}

func ProcessRandom(text string, words *[]string) (string, error) {
  return getDefaultProcessor().ProcessRandom(text, words)
}

func ProcessExt(text string, raw_variants []string, variants [][]string, indexes []int, words *[]string) (string, bool, error) {
  var vr string
  for i, v := range raw_variants {
    if nil != indexes && len(indexes) > i {
      vr = variants[i][indexes[i]]
    } else {
      vr = variant(variants[i])
    }
    if nil != words {
      *words = append(*words, vr)
    }
    text = strings.Replace(text, v, vr, 1)
  }
  return text, incIndex(variants, indexes) /* can continue? */, nil
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Generators
///////////////////////////////////////////////////////////////////////////////

func (p *Processor) GenerateRandom(text string, count, try_count int, exclude bool) chan string {
  var excludes []string = nil
  if count < 1 {
    count = -1
    exclude = true
  }
  if try_count < 1 {
    try_count = 1
  }

  if exclude {
    excludes = make([]string, 0)
  }

  ch := make(chan string)

  go func() {
    raw_variants, variants, err := p.PrepareVariants(text)
    if nil == variants || len(variants) < 1 {
      err = errorNoVariants
    }

    if nil == err {
      for {
        if count == 0 {
          break
        }
        if count > 0 {
          count--
        }

        b := false
        for i := 0; i < try_count; i++ {
          s, _, e := ProcessExt(text, raw_variants, variants, nil, nil)
          if nil == e {
            if !containt(excludes, s) {
              if nil != excludes {
                excludes = append(excludes, s)
              }
              b = true
              ch <- s
              break
            }
          }
        } // end for

        if !b {
          break
        }
      }
    }
    close(ch)
  }()

  return ch
}

func GenerateRandom(text string, count, try_count int, exclude bool) chan string {
  return getDefaultProcessor().GenerateRandom(text, count, try_count, exclude)
}

func (p *Processor) Generate(text string, count int) chan string {
  if count < 1 {
    count = -1
  }

  ch := make(chan string)

  go func() {
    raw_variants, variants, err := p.PrepareVariants(text)
    if nil == variants || len(variants) < 1 {
      err = errorNoVariants
    }

    indexes := make([]int, len(raw_variants))

    if nil == err {
      for {
        if count == 0 {
          break
        }
        if count > 0 {
          count--
        }

        s, b, e := ProcessExt(text, raw_variants, variants, indexes, nil)
        if nil == e {
          ch <- s
        }
        if !b {
          break
        }
      }
    }
    close(ch)
  }()

  return ch
}

func Generate(text string, count int) chan string {
  return getDefaultProcessor().Generate(text, count)
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Helpers
///////////////////////////////////////////////////////////////////////////////

func variant(variants []string) string {
  return variants[randInt(0, len(variants)-1)]
}

func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}

func containt(arr []string, s string) bool {
  if nil != arr {
    for _, it := range arr {
      if it == s {
        return true
      }
    }
  }
  return false
}

func incIndex(variants [][]string, indexes []int) bool {
  if nil == variants || nil == indexes || len(indexes) < len(variants) {
    return true
  }
  last := len(variants) - 1
  return incIndexBy(last, variants, indexes)
}

func incIndexBy(i int, variants [][]string, indexes []int) bool {
  indexes[i]++
  if len(variants[i]) <= indexes[i] {
    if i < 1 {
      return false
    }
    indexes[i] = 0
    if i > 0 {
      return incIndexBy(i-1, variants, indexes)
    }
  }
  return true
}
