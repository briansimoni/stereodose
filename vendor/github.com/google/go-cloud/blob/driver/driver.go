// Copyright 2018 The Go Cloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package driver defines a set of interfaces that the blob package uses to interact
// with the underlying blob services.
package driver

import (
	"context"
	"io"
	"time"
)

// ErrorKind is a code to indicate the kind of failure.
type ErrorKind int

const (
	GenericError ErrorKind = iota
	NotFound
)

// Error is an interface that may be implemented by an error returned by
// a driver to indicate the kind of failure.  If an error does not have the
// Kind method, then it is assumed to be GenericError.
type Error interface {
	error
	Kind() ErrorKind
}

// Reader reads an object from the blob.
type Reader interface {
	io.ReadCloser

	// Attributes returns a subset of attributes about the blob.
	// Use Bucket.Attributes to get the full set.
	Attributes() ReaderAttributes
}

// Writer writes an object to the blob.
type Writer interface {
	io.WriteCloser
}

// WriterOptions controls behaviors of Writer.
type WriterOptions struct {
	// BufferSize changes the default size in byte of the maximum part Writer can
	// write in a single request, if supported. Larger objects will be split into
	// multiple requests.
	BufferSize int
	// Metadata holds key/value strings to be associated with the blob.
	// Keys are guaranteed to be non-empty and lowercased.
	Metadata map[string]string
}

// ReaderAttributes contains a subset of attributes about a blob that are
// accessible from Reader. Use Bucket.Attributes to get the full set of
// attributes.
type ReaderAttributes struct {
	// ContentType is the MIME type of the blob object. It must not be empty.
	ContentType string
	// ModTime is the time the blob object was last modified.
	ModTime time.Time
	// Size is the size of the object in bytes.
	Size int64
}

// Attributes contains attributes about a blob.
type Attributes struct {
	// ContentType is the MIME type of the blob object. It must not be empty.
	ContentType string
	// Metadata holds key/value pairs associated with the blob.
	// Keys will be lowercased by the concrete type before being returned
	// to the user. If there are duplicate case-insensitive keys (e.g.,
	// "foo" and "FOO"), only one value will be kept, and it is undefined
	// which one.
	Metadata map[string]string
	// ModTime is the time the blob object was last modified.
	ModTime time.Time
	// Size is the size of the object in bytes.
	Size int64
}

// Bucket provides read, write and delete operations on objects within it on the
// blob service.
type Bucket interface {
	// Attributes returns attributes for the blob. If the specified object does
	// not exist, Attributes must return an error whose Kind method returns
	// NotFound.
	Attributes(ctx context.Context, key string) (Attributes, error)

	// NewRangeReader returns a Reader that reads part of an object, reading at
	// most length bytes starting at the given offset. If length is negative, it
	// will read till the end of the object. If the specified object does not
	// exist, NewRangeReader must return an error whose Kind method returns
	// NotFound.
	NewRangeReader(ctx context.Context, key string, offset, length int64) (Reader, error)

	// NewTypedWriter returns Writer that writes to an object associated with key.
	//
	// A new object will be created unless an object with this key already exists.
	// Otherwise any previous object with the same name will be replaced.
	// The object may not be available (and any previous object will remain)
	// until Close has been called.
	//
	// contentType sets the MIME type of the object to be written. It must not be
	// empty.
	//
	// The caller must call Close on the returned Writer when done writing.
	//
	// Implementations should abort an ongoing write if ctx is later canceled,
	// and do any necessary cleanup in Close. Close should then return ctx.Err().
	NewTypedWriter(ctx context.Context, key string, contentType string, opt *WriterOptions) (Writer, error)

	// Delete deletes the object associated with key. If the specified object does
	// not exist, NewRangeReader must return an error whose Kind method
	// returns NotFound.
	Delete(ctx context.Context, key string) error
}
