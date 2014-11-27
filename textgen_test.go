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
  "math/rand"
  "testing"
  "time"
)

var (
  text = "[Test|Text|Guest|Start it|Crash it] [for|not for] [example|fun|you|us]! Oh ya!"
)

func init() {
  rand.Seed(time.Now().UTC().UnixNano() ^ int64(time.Now().Nanosecond()))
}

///////////////////////////////////////////////////////////////////////////////
/// MARK: Tests
///////////////////////////////////////////////////////////////////////////////

func TestPrepareVariants(t *testing.T) {
  all, variants, err := PrepareVariants(text)
  max := 1

  if nil != variants {
    for _, a := range variants {
      max *= len(a)
    }
  }

  t.Logf("All Variants: %v", all)
  t.Logf("Variants: %v, %v", variants, err)
  t.Logf("Variants MAX: %d", max)

  if nil != err {
    t.Error(err)
  }
}

func TestProcessRandom(t *testing.T) {
  gen_text, err := ProcessRandom(text, nil)
  t.Logf("ProcessRandom: %v, %v", gen_text, err)

  if nil != err {
    t.Error(err)
  }
}

func TestGenerateRandom(t *testing.T) {
  i := 1
  for s := range GenerateRandom(text, 0, 10, true) {
    t.Logf("GenerateRandom: %d) %v", i, s)
    i++
  }
}

func TestGenerate(t *testing.T) {
  i := 1
  for s := range Generate(text, 0) {
    t.Logf("Generate: %d) %v", i, s)
    i++
  }
}
