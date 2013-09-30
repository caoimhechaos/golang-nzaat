// Copyright 2013 Caoimhe Chaos <caoimhechaos@protonmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

package nzaat

import (
	"hash"
	"testing"
)

// Test the hash of an empty string.
func TestEmpty(t *testing.T) {
	var h hash.Hash32 = New()
	var res uint32 = h.Sum32()

	t.Logf("NZAAT(\"\") = %x\n", res)

	if res != 0 {
		t.Fail()
	}
}

// Test the hash of a string with just "a" in it.
func TestStringA(t *testing.T) {
	var h hash.Hash32 = New()
	var res uint32

	h.Write([]byte("a"))
	res = h.Sum32()

	t.Logf("NZAAT(\"a\") = %x\n", res)

	if res != 0xc31517c4 {
		t.Fail()
	}
}

// Test the hash of a string with "abc" in it.
func TestStringABC(t *testing.T) {
	var h hash.Hash32 = New()
	var res uint32

	h.Write([]byte("abc"))
	res = h.Sum32()

	t.Logf("NZAAT(\"abc\") = %x\n", res)

	if res != 0xC3E39E2D {
		t.Fail()
	}
}

// Test the hash of a string with "message digest" in it.
func TestStringMessageDigest(t *testing.T) {
	var h hash.Hash32 = New()
	var res uint32

	h.Write([]byte("message digest"))
	res = h.Sum32()

	t.Logf("NZAAT(\"message digest\") = %x\n", res)

	if res != 0x434B78B4 {
		t.Fail()
	}
}
