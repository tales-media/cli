/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Body interface {
	Len() int64
	Reader() (io.ReadCloser, error)
	ContentType() string
}

var NoBody = noBody{}

type noBody struct{}

var _ Body = noBody{}

func (noBody) Len() int64                     { return 0 }
func (noBody) Reader() (io.ReadCloser, error) { return http.NoBody, nil }
func (noBody) ContentType() string            { return "" }

type bytesBody struct {
	b  []byte
	ct string
}

var _ Body = &bytesBody{}

func NewBody(b []byte, contentType string) *bytesBody {
	return &bytesBody{
		b:  b,
		ct: contentType,
	}
}

func NewBufferBody(buf bytes.Buffer, contentType string) *bytesBody {
	return &bytesBody{
		b:  buf.Bytes(),
		ct: contentType,
	}
}

func (b *bytesBody) Len() int64 {
	return int64(len(b.b))
}

func (b *bytesBody) Reader() (io.ReadCloser, error) {
	r := bytes.NewReader(b.b)
	return io.NopCloser(r), nil
}

func (b *bytesBody) ContentType() string {
	return b.ct
}

type stringReaderBody struct {
	sr *strings.Reader
	ct string
}

var _ Body = &stringReaderBody{}

func NewStringBody(s string, contentType string) *stringReaderBody {
	return &stringReaderBody{
		sr: strings.NewReader(s),
		ct: contentType,
	}
}

func NewStringReaderBody(sr *strings.Reader, contentType string) *stringReaderBody {
	return &stringReaderBody{
		sr: sr,
		ct: contentType,
	}
}

func (b *stringReaderBody) Len() int64 {
	return int64(b.sr.Len())
}

func (b *stringReaderBody) Reader() (io.ReadCloser, error) {
	copy := *b.sr
	return io.NopCloser(&copy), nil
}

func (b *stringReaderBody) ContentType() string {
	return b.ct
}

type EncoderFunc func(v any) ([]byte, error)

type encoderBody struct {
	v    any
	f    EncoderFunc
	b    []byte
	err  error
	once sync.Once
	ct   string
}

var _ Body = &encoderBody{}

func NewEncoderBody(v any, contentType string, f EncoderFunc) *encoderBody {
	return &encoderBody{
		v:  v,
		f:  f,
		ct: contentType,
	}
}

func (b *encoderBody) Len() int64 {
	_ = b.Encode()
	return int64(len(b.b))
}

func (b *encoderBody) Reader() (io.ReadCloser, error) {
	if err := b.Encode(); err != nil {
		return nil, err
	}
	r := bytes.NewReader(b.b)
	return io.NopCloser(r), nil
}

func (b *encoderBody) Encode() error {
	b.once.Do(func() {
		b.b, b.err = b.f(b.v)
	})
	return b.err
}

func (b *encoderBody) ContentType() string {
	return b.ct
}
