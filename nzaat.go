// Copyright 2013 Tonnerre Lombard <tonnerre@ancient-solutions.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// NZAAT originates from Thorsten “mirabilos” Glaser <tg@mirbsd.org>.
// This file defines the NZAAT hash which is derived from Bob Jenkins’
// one at a time hash with the goals to ① not have 00000000 as result
// (for speed optimisations with hash table lookups) and ② change for
// every input octet, even NUL bytes – no matter what state. For con‐
// venience, NZAAT, which can have an all-zero result, for outside of
// hash table lookup dereference deferral, is also provided.
//
// We define the following primitives (all modulo 2³²):
//
// OUP(s,b) → { OAV(s,b); MIX(s); }
//      This is the Update function of Jenkins’ one-at-a-time hash
// NUP(s,b) → { NAV(s,b); MIX(s); }
//      This shall be the body of the for loop in the Write method
// OAV(s,b) → { ADD(s,b); }
//      This is how the Add part of Jenkins’ one-at-a-time hash works
// NAV(s,b) → { ADD(s,b+c); } – c const
//      This is how the Add part of Write differs – by a const
// ADD(x,y) → { x += y; }
//      Generic addition function
// MIX(s) → { s += s << 10; s ^= s >> 6; }
//      This is the Mix part of both Update functions (common)
// FIN(s) → { s += s << 3; s ^= s >> 11; s += s << 15; result ← s; }
//      This is the Postprocess function for the one-at-a-time hash
// NZF(s) → { MIX(s); if(!s){++s;} FIN(s); }
// NZF(s) ≘ { if (!s) { s = 1; } else { NAF(s); } }
//      This is the Postprocess function for the new NZAT hash
// NAF(s) → { MIX(s); FIN(s); }
//      This shall become the Postprocess function for the NZAAT hash
//
// This means that the difference between OAAT and NZAT is ① a factor
// c (constant) that is added to the state in addition to every input
// octet, ② preventing the hash result from becoming 0, by a volunta‐
// ry collision on exactly(!) one value, and ③ mixing in another zero
// data octet to improve avalanche behaviour; OAAT’s avalanche bitmap
// is all green except for the last (OAAT) data octet; we avoid a few
// yellow-bit patterns, nothing really bad, by making that last octet
// be not really the last. The impact of these changes doesn’t worsen
// the hash’s properties, except for NZAT’s new collision. NZAAT does
// not collide more than OAAT, so it’s better in all respects.
//
// First some observations on some of the above primitives, proven by
// me empirically (can be proven algebraically as well):
//
// FIN(s) = 0 ⇐ s = 0
//      Only a state of 0 yields a hash result of 0
// MIX(s) = 0 ⇐ s = 0
//      When updating only the Add part (not the Mix part) determines
//      whether the next state is 0
// ADD(s,v) = 0 ⇐ s = 100000000h - v (mod 2³²)
//      For each input byte value there is exactly one previous-state
//      producing a 0 next-state (that being -s mod 2³² for dword in‐
//      put values and, for bytes, simpler -s mod 2⁸)
// ADD(s,b+c) ≠ s ⇐ b is byte, c is 1, works
//      Choose c arbitrarily, 1 is economical, even on e.g. ARM
//
// All this means we can use an IV of 0, add arbitrary octet data and
// have a good result (for NZAAT) while having to provide for exactly
// one 2-in-1 collision for NZAT manually. This implementation mainly
// aims for economic code (small, fast, RISC CPUs notwithstanding).

package nzaat

import "hash"

type digest uint32

// New returns a new hash.Hash32 computing the NZAAT checksum.
func New() hash.Hash32 {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Reset() {
	*d = 0
}

func (d *digest) Size() int {
	return 4
}

func (d *digest) BlockSize() int {
	return 1
}

func (d *digest) Write(p []byte) (nn int, err error) {
	for _, x := range p {
		*d += digest(x) + 1
		*d += *d << 10
		*d ^= *d >> 6
	}

	return len(p), nil
}

// Count NILs in all parts
func (d *digest) Sum32() uint32 {
	var sum uint32 = uint32(*d)

	sum += sum << 10
	sum ^= sum >> 6
	sum += sum << 3
	sum ^= sum >> 11
	sum += sum << 15

	return sum
}

func (d *digest) Sum(in []byte) []byte {
	var s uint32 = d.Sum32()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

// Checksum returns the NZAAT checksum of data.
func Checksum(data []byte) uint32 {
	var h hash.Hash32 = New()
	h.Write(data)
	return h.Sum32()
}
