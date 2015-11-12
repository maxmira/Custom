package custom

import (
 "github.com/AlasdairF/Unleak"
 "unicode/utf8"
 "io"
 "errors"
)

const (
	bestLength  = 100000 // determined in trials on writing to file and writing to disk
	bestLengthMinus1  = bestLength - 1
	bestLengthMinus2  = bestLength - 2
	bestLengthMinus3  = bestLength - 3
	bestLengthMinus4  = bestLength - 4
	bestLengthMinus5  = bestLength - 5
	bestLengthMinus6  = bestLength - 6
	bestLengthMinus7  = bestLength - 7
	bestLengthMinus8  = bestLength - 8
)

// Constants stolen from unicode/utf8 for WriteRune
const (
	maxRune   = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000
	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111
	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

// -------- FIXED BUFFER WRITER --------

type Writer struct {
	w io.Writer
	data [bestLength]byte
	cursor int
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (b *Writer) Write(p []byte) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.w.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Writer
			b.cursor = 0
			return b.w.Write(p)
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Writer) WriteString(p string) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.w.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Writer
			b.cursor = 0
			return b.w.Write([]byte(p))
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Writer) WriteByte(p byte) error {
	if b.cursor < bestLength {
		b.data[b.cursor] = p
		b.cursor++
		return nil
	}
	var err error
	if b.cursor > 0 {
		_, err = b.w.Write(b.data[0:b.cursor]) // flush
	}
	b.data[0] = p
	b.cursor = 1
	return err
}

func (b *Writer) WriteRune(r rune) (int, error) {
	switch i := uint32(r); {
	case i <= rune1Max:
		err := b.WriteByte(byte(r))
		return 1, err
	case i <= rune2Max:
		err := b.Write2Bytes(t2 | byte(r>>6), tx | byte(r)&maskx)
		return 2, err
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = '\uFFFD'
		fallthrough
	case i <= rune3Max:
		err := b.Write3Bytes(t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 3, err
	default:
		err := b.Write4Bytes(t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 4, err
	}
}

func (b *Writer) Write2Bytes(p1, p2 byte) error {
	if b.cursor < bestLengthMinus1 {
		b.data[b.cursor] = p1
		b.data[b.cursor + 1] = p2
		b.cursor += 2
		return nil
	}
	var err error
	if b.cursor > 0 {
		_, err = b.w.Write(b.data[0:b.cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.cursor = 2
	return err
}

func (b *Writer) Write3Bytes(p1, p2, p3 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus2 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.cursor += 3
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.cursor = 3
	return err
}

func (b *Writer) Write4Bytes(p1, p2, p3, p4 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus3 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.cursor += 4
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.cursor = 4
	return err
}

func (b *Writer) Write5Bytes(p1, p2, p3, p4, p5 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus4 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.data[cursor + 4] = p5
		b.cursor += 5
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.data[4] = p5
	b.cursor = 5
	return err
}

func (b *Writer) Write6Bytes(p1, p2, p3, p4, p5, p6 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus5 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.data[cursor + 4] = p5
		b.data[cursor + 5] = p6
		b.cursor += 6
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.data[4] = p5
	b.data[5] = p6
	b.cursor = 6
	return err
}

func (b *Writer) Write7Bytes(p1, p2, p3, p4, p5, p6, p7 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus6 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.data[cursor + 4] = p5
		b.data[cursor + 5] = p6
		b.data[cursor + 6] = p7
		b.cursor += 7
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.data[4] = p5
	b.data[5] = p6
	b.data[6] = p7
	b.cursor = 7
	return err
}

func (b *Writer) Write8Bytes(p1, p2, p3, p4, p5, p6, p7, p8 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus7 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.data[cursor + 4] = p5
		b.data[cursor + 5] = p6
		b.data[cursor + 6] = p7
		b.data[cursor + 7] = p8
		b.cursor += 8
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.data[4] = p5
	b.data[5] = p6
	b.data[6] = p7
	b.data[7] = p8
	b.cursor = 8
	return err
}

func (b *Writer) Write9Bytes(p1, p2, p3, p4, p5, p6, p7, p8, p9 byte) error {
	cursor := b.cursor
	if cursor < bestLengthMinus8 {
		b.data[cursor] = p1
		b.data[cursor + 1] = p2
		b.data[cursor + 2] = p3
		b.data[cursor + 3] = p4
		b.data[cursor + 4] = p5
		b.data[cursor + 5] = p6
		b.data[cursor + 6] = p7
		b.data[cursor + 7] = p8
		b.data[cursor + 8] = p9
		b.cursor += 9
		return nil
	}
	var err error
	if cursor > 0 {
		_, err = b.w.Write(b.data[0:cursor]) // flush
	}
	b.data[0] = p1
	b.data[1] = p2
	b.data[2] = p3
	b.data[3] = p4
	b.data[4] = p5
	b.data[5] = p6
	b.data[6] = p7
	b.data[7] = p8
	b.data[8] = p9
	b.cursor = 9
	return err
}

func (w *Writer) WriteBool(v bool) {
	if v {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (w *Writer) Write2Bools(v1, v2 bool) {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	w.WriteByte(b)
}

func (w *Writer) Write2Unit4s(v1, v2 uint8) {
	v1 |= v2 << 4
	w.WriteByte(v1)
}

func (w *Writer) WriteUint16(v uint16) {
	w.Write2Bytes(byte(v), byte(v >> 8))
}

// If it's less than 255 then it's encoded in the 1st byte, otherwise 1st byte is 255 and it's encoded in two more bytes
// This is only useful if it is expected that the value will be <255 more than half the time
func (w *Writer) WriteUint16Variable(v uint16) {
	if v < 255 {
		w.WriteByte(byte(v))
		return
	}
	w.Write3Bytes(255, byte(v), byte(v >> 8))
}

func (w *Writer) WriteInt16Variable(v int16) {
	if v > -128 && v < 128 {
		w.WriteByte(byte(v + 127))
		return
	}
	v2 := uint16(v)
	w.Write3Bytes(255, byte(v2), byte(v2 >> 8))
}

func (w *Writer) WriteUint24(v uint32) {
	w.Write3Bytes(byte(v), byte(v >> 8), byte(v >> 16))
}

func (w *Writer) WriteUint32(v uint32) {
	w.Write4Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
}

func (w *Writer) WriteUint48(v uint64) {
	w.Write6Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
}

func (w *Writer) WriteUint64(v uint64) {
	w.Write8Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (w *Writer) WriteUint64Variable(v uint64) {
	switch numbytes(v) {
		case 0: w.WriteByte(0)
		case 1: w.Write2Bytes(1, byte(v))
		case 2: w.Write3Bytes(2, byte(v), byte(v >> 8))
		case 3: w.Write4Bytes(3, byte(v), byte(v >> 8), byte(v >> 16))
		case 4: w.Write5Bytes(4, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
		case 5: w.Write6Bytes(5, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32))
		case 6: w.Write7Bytes(6, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
		case 7: w.Write8Bytes(7, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48))
		case 8: w.Write9Bytes(8, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 25), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
	}
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (w *Writer) Write2Uint64sVariable(v1 uint64, v2 uint64) {
	s1 := numbytes(v1)
	s2 := numbytes(v2)
	w.WriteByte((s1 << 4) | s2)
	switch s1 {
		case 1: w.WriteByte(byte(v1))
		case 2: w.Write2Bytes(byte(v1), byte(v1 >> 8))
		case 3: w.Write3Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16))
		case 4: w.Write4Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24))
		case 5: w.Write5Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32))
		case 6: w.Write6Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40))
		case 7: w.Write7Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48))
		case 8: w.Write8Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 25), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48), byte(v1 >> 56))
	}
	switch s2 {
		case 1: w.WriteByte(byte(v2))
		case 2: w.Write2Bytes(byte(v2), byte(v2 >> 8))
		case 3: w.Write3Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16))
		case 4: w.Write4Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24))
		case 5: w.Write5Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32))
		case 6: w.Write6Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40))
		case 7: w.Write7Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48))
		case 8: w.Write8Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 25), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48), byte(v2 >> 56))
	}
}

func (w *Writer) WriteFloat32(flt float32) {
	w.WriteUint32(math.Float32bits(flt))
}

func (w *Writer) WriteFloat64(flt float64) {
	w.WriteUint64(math.Float64bits(flt))
}

func (w *Writer) WriteString8(s string) {
	l := len(s)
	w.WriteByte(uint8(l))
	if l > 255 {
		w.WriteString(s[0:255])
	} else {
		w.WriteString(s)
	}
}

func (w *Writer) WriteString16(s string) {
	l := len(s)
	w.WriteUint16(uint16(l))
	if l > 65535 {
		w.WriteString(s[0:65535])
	} else {
		w.WriteString(s)
	}
}

func (w *Writer) WriteString32(s string) {
	l := len(s)
	w.WriteUint32(uint32(l))
	if l > 4294967295 {
		w.WriteString(s[0:4294967295])
	} else {
		w.WriteString(s)
	}
}

/*
// 12 bits and 4 bits
func (w *Writer) Write12And4(v1, v2 uint16) {
	v1 |= v2 << 12
	w.Write2Bytes(byte(v1), byte(v1 >> 8))
}
*/

func (w *Writer) WriteSpecial(v1 uint8, b1, b2, b3, b4 bool) {
	if b1 {
		v1 |= 128
	}
	if b2 {
		v1 |= 64
	}
	if b3 {
		v1 |= 32
	}
	if b4 {
		v1 |= 16
	}
	w.WriteByte(v1)
}

func (w *Writer) WriteSpecial2(v1, v2, v3 uint8, b1 bool) {
	v1 |= v2 << 3
	v1 |= v3 << 5
	if b1 {
		v1 |= 128
	}
	w.WriteByte(v1)
}

func (b *Writer) Close() (err error) {
	if b.cursor > 0 {
		_, err = b.w.Write(b.data[0:b.cursor])
		b.cursor = 0
	}
	if sw, ok := b.w.(io.Closer); ok { // Attempt to close underlying writer if it has a Close() method
		if err == nil {
			err = sw.Close()
		} else {
			sw.Close()
		}
	}
	b.w = nil
	return
}

func (b *Writer) Flush() (err error) {
	if b.cursor > 0 {
		_, err = b.w.Write(b.data[0:b.cursor])
		b.cursor = 0
	}
	return
}

func (b *Writer) Recycle(w io.Writer) (err error) {
	b.cursor = 0
	b.w = w
	return
}

// -------- GROWING BUFFER --------

type Buffer struct {
	data []byte
	cursor, length int
}

func NewBuffer(l int) *Buffer {
	return &Buffer{data: make([]byte, l), length: l}
}

func (b *Buffer) Write(p []byte) (int, error) {
	l := len(p)
	if b.cursor + l > b.length {
		b.length = (b.length + l) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Buffer) WriteString(p string) (int, error) {
	l := len(p)
	if b.cursor + l > b.length {
		b.length = (b.length + l) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Buffer) WriteByte(p byte) error {
	if b.cursor >= b.length {
		b.length = (b.length + 1) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[b.cursor] = p
	b.cursor++
	return nil
}

func (b *Buffer) WriteRune(r rune) (int, error) {
	switch i := uint32(r); {
	case i <= rune1Max:
		err := b.WriteByte(byte(r))
		return 1, err
	case i <= rune2Max:
		err := b.Write2Bytes(t2 | byte(r>>6), tx | byte(r)&maskx)
		return 2, err
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = '\uFFFD'
		fallthrough
	case i <= rune3Max:
		err := b.Write3Bytes(t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 3, err
	default:
		err := b.Write4Bytes(t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx)
		return 4, err
	}
}

func (b *Buffer) Write2Bytes(p1, p2 byte) error {
	c := b.cursor
	if c + 2 > b.length {
		b.length = (b.length + 2) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.cursor += 2
	return nil
}

func (b *Buffer) Write3Bytes(p1, p2, p3 byte) error {
	c := b.cursor
	if c + 3 > b.length {
		b.length = (b.length + 3) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.cursor += 3
	return nil
}

func (b *Buffer) Write4Bytes(p1, p2, p3, p4 byte) error {
	c := b.cursor
	if c + 4 > b.length {
		b.length = (b.length + 4) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.cursor += 4
	return nil
}

func (b *Buffer) Write5Bytes(p1, p2, p3, p4, p5 byte) error {
	c := b.cursor
	if c + 5 > b.length {
		b.length = (b.length + 5) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.data[c + 4] = p5
	b.cursor += 5
	return nil
}

func (b *Buffer) Write6Bytes(p1, p2, p3, p4, p5, p6 byte) error {
	c := b.cursor
	if c + 6 > b.length {
		b.length = (b.length + 6) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.data[c + 4] = p5
	b.data[c + 5] = p6
	b.cursor += 6
	return nil
}

func (b *Buffer) Write7Bytes(p1, p2, p3, p4, p5, p6, p7 byte) error {
	c := b.cursor
	if c + 7 > b.length {
		b.length = (b.length + 7) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.data[c + 4] = p5
	b.data[c + 5] = p6
	b.data[c + 6] = p7
	b.cursor += 7
	return nil
}

func (b *Buffer) Write8Bytes(p1, p2, p3, p4, p5, p6, p7, p8 byte) error {
	c := b.cursor
	if c + 8 > b.length {
		b.length = (b.length + 8) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.data[c + 4] = p5
	b.data[c + 5] = p6
	b.data[c + 6] = p7
	b.data[c + 7] = p8
	b.cursor += 8
	return nil
}

func (b *Buffer) Write9Bytes(p1, p2, p3, p4, p5, p6, p7, p8, p9 byte) error {
	c := b.cursor
	if c + 9 > b.length {
		b.length = (b.length + 9) * 2
		newAr := make([]byte, b.length)
		copy(newAr, b.data)
		b.data = newAr
	}
	b.data[c] = p1
	b.data[c + 1] = p2
	b.data[c + 2] = p3
	b.data[c + 3] = p4
	b.data[c + 4] = p5
	b.data[c + 5] = p6
	b.data[c + 6] = p7
	b.data[c + 7] = p8
	b.data[c + 8] = p9
	b.cursor += 9
	return nil
}

func (b *Buffer) WriteBool(v bool) {
	if v {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (b *Buffer) Write2Bools(v1, v2 bool) {
	var b byte
	if v1 {
		b = 1
	}
	if v2 {
		b |= 2
	}
	b.WriteByte(b)
}

func (b *Buffer) Write2Unit4s(v1, v2 uint8) {
	v1 |= v2 << 4
	b.WriteByte(v1)
}

func (b *Buffer) WriteUint16(v uint16) {
	b.Write2Bytes(byte(v), byte(v >> 8))
}

// If it's less than 255 then it's encoded in the 1st byte, otherwise 1st byte is 255 and it's encoded in two more bytes
// This is only useful if it is expected that the value will be <255 more than half the time
func (b *Buffer) WriteUint16Variable(v uint16) {
	if v < 255 {
		b.WriteByte(byte(v))
		return
	}
	b.Write3Bytes(255, byte(v), byte(v >> 8))
}

func (b *Buffer) WriteInt16Variable(v int16) {
	if v > -128 && v < 128 {
		b.WriteByte(byte(v + 127))
		return
	}
	v2 := uint16(v)
	b.Write3Bytes(255, byte(v2), byte(v2 >> 8))
}

func (b *Buffer) WriteUint24(v uint32) {
	b.Write3Bytes(byte(v), byte(v >> 8), byte(v >> 16))
}

func (b *Buffer) WriteUint32(v uint32) {
	b.Write4Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
}

func (b *Buffer) WriteUint48(v uint64) {
	b.Write6Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
}

func (b *Buffer) WriteUint64(v uint64) {
	b.Write8Bytes(byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (b *Buffer) WriteUint64Variable(v uint64) {
	switch numbytes(v) {
		case 0: b.WriteByte(0)
		case 1: b.Write2Bytes(1, byte(v))
		case 2: b.Write3Bytes(2, byte(v), byte(v >> 8))
		case 3: b.Write4Bytes(3, byte(v), byte(v >> 8), byte(v >> 16))
		case 4: b.Write5Bytes(4, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24))
		case 5: b.Write6Bytes(5, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32))
		case 6: b.Write7Bytes(6, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40))
		case 7: b.Write8Bytes(7, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48))
		case 8: b.Write9Bytes(8, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 25), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56))
	}
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (b *Buffer) Write2Uint64sVariable(v1 uint64, v2 uint64) {
	s1 := numbytes(v1)
	s2 := numbytes(v2)
	b.WriteByte((s1 << 4) | s2)
	switch s1 {
		case 1: b.WriteByte(byte(v1))
		case 2: b.Write2Bytes(byte(v1), byte(v1 >> 8))
		case 3: b.Write3Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16))
		case 4: b.Write4Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24))
		case 5: b.Write5Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32))
		case 6: b.Write6Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40))
		case 7: b.Write7Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 24), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48))
		case 8: b.Write8Bytes(byte(v1), byte(v1 >> 8), byte(v1 >> 16), byte(v1 >> 25), byte(v1 >> 32), byte(v1 >> 40), byte(v1 >> 48), byte(v1 >> 56))
	}
	switch s2 {
		case 1: b.WriteByte(byte(v2))
		case 2: b.Write2Bytes(byte(v2), byte(v2 >> 8))
		case 3: b.Write3Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16))
		case 4: b.Write4Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24))
		case 5: b.Write5Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32))
		case 6: b.Write6Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40))
		case 7: b.Write7Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 24), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48))
		case 8: b.Write8Bytes(byte(v2), byte(v2 >> 8), byte(v2 >> 16), byte(v2 >> 25), byte(v2 >> 32), byte(v2 >> 40), byte(v2 >> 48), byte(v2 >> 56))
	}
}

func (b *Buffer) WriteFloat32(flt float32) {
	b.WriteUint32(math.Float32bits(flt))
}

func (b *Buffer) WriteFloat64(flt float64) {
	b.WriteUint64(math.Float64bits(flt))
}

func (b *Buffer) WriteString8(s string) {
	l := len(s)
	b.WriteByte(uint8(l))
	if l > 255 {
		b.WriteString(s[0:255])
	} else {
		b.WriteString(s)
	}
}

func (b *Buffer) WriteString16(s string) {
	l := len(s)
	b.WriteUint16(uint16(l))
	if l > 65535 {
		b.WriteString(s[0:65535])
	} else {
		b.WriteString(s)
	}
}

func (b *Buffer) WriteString32(s string) {
	l := len(s)
	b.WriteUint32(uint32(l))
	if l > 4294967295 {
		b.WriteString(s[0:4294967295])
	} else {
		b.WriteString(s)
	}
}

/*
// 12 bits and 4 bits
func (b *Buffer) Write12And4(v1, v2 uint16) {
	v1 |= v2 << 12
	b.Write2Bytes(byte(v1), byte(v1 >> 8))
}
*/

func (b *Buffer) WriteSpecial(v1 uint8, b1, b2, b3, b4 bool) {
	if b1 {
		v1 |= 128
	}
	if b2 {
		v1 |= 64
	}
	if b3 {
		v1 |= 32
	}
	if b4 {
		v1 |= 16
	}
	b.WriteByte(v1)
}

func (b *Buffer) WriteSpecial2(v1, v2, v3 uint8, b1 bool) {
	v1 |= v2 << 3
	v1 |= v3 << 5
	if b1 {
		v1 |= 128
	}
	b.WriteByte(v1)
}

func (b *Buffer) Reset() {
	b.cursor = 0
	return
}

func (b *Buffer) Len() int {
	return b.cursor
}

func (b *Buffer) Bytes() []byte {
	return b.data[0:b.cursor]
}

func (b *Buffer) String() string {
	return string(b.data[0:b.cursor])
}

func numbytes(v uint64) uint8 {
	switch {
		case v == 0: return 0
		case v < 256: return 1
		case v < 65536: return 2
		case v < 16777216: return 3
		case v < 4294967296: return 4
		case v < 1099511627776: return 5
		case v < 281474976710655: return 6
		case v < 72057594037927936: return 7
		default: return 8
	}
}

// -------- READER --------

type Reader struct {
	f io.Reader
	at int		// the cursor for where I am in buf
	n int		// how much uncompressed but as of yet unparsed data is left in buf
	buf []byte	// the buffer for reading data
}

func NewReader(f io.Reader, bufsize int) *Reader {
	return &Reader{f: f, buf: make([]byte, bufsize + 512)} // 512 is bytes.MinRead
}

func (r *Reader) Read(b []byte) {
	x := len(b)
	for r.n < x {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	copy(b, r.buf[r.at:r.at+x]) // must be copied to avoid memory leak
	r.at += x
	r.n -= x
}

func (r *Reader) Readx(x int) []byte {
	for r.n < x {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	tmp := make([]byte, x)
	copy(tmp, r.buf[r.at:r.at+x]) // must be copied to avoid memory leak
	r.at += x
	r.n -= x
	return tmp
}

func (r *Reader) ReadBool() bool {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	var b1 bool
	if r.buf[r.at] > 0 {
		b1 = true
	}
	r.at++
	r.n--
	return b1
}

func (r *Reader) Read2Bools() (bool, bool) {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	var b1, b2 bool
	switch r.buf[r.at] {
		case 1: b1 = true
		case 2: b2 = true
		case 3: b1, b2 = true, true
	}
	r.at++
	r.n--
	return b1, b2
}

func (r *Reader) Read2Unit4s() (uint8, uint8) {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	res1, res2 := r.buf[r.at] & 15, r.buf[r.at] >> 4
	r.at++
	r.n--
	return res1, res2
}

func (r *Reader) ReadByte() uint8 {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	res := r.buf[r.at]
	r.at++
	r.n--
	return res
}

func (r *Reader) ReadUTF8() []byte {
	first := r.ReadByte()
	if first < 128 { // length 1
		return []byte{first}
	}
	if first & 32 == 0 { // length 2
			return []byte{first, r.ReadByte()}
	} else {
		b := make([]byte, 3)
		b[0] = first
		b[1] = r.ReadByte()
		b[2] = r.ReadByte()
		return b
	}
}

func (r *Reader) ReadRune() rune {
	first := r.ReadByte()
	if first < 128 { // length 1
		return rune(first)
	}
	if first & 32 == 0 { // length 2
		r, _ := utf8.DecodeRune([]byte{first, r.ReadByte()})
		return r
	} else {
		b := make([]byte, 3)
		b[0] = first
		b[1] = r.ReadByte()
		b[2] = r.ReadByte()
		r, _ := utf8.DecodeRune(b)
		return r
	}
}

func (r *Reader) ReadUint16() uint16 {
	for r.n < 2 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint16(r.buf[r.at]) | uint16(r.buf[r.at+1])<<8
	r.at += 2
	r.n -= 2
	return res
}

// If it's less than 255 then it's encoded in the 1st byte, otherwise 1st byte is 255 and it's encoded in two more bytes
// This is only useful if it is expected that the value will be <255 more than half the time
func (r *Reader) ReadUint16Variable() uint16 {
	v := r.ReadByte()
	if v < 255 {
		return uint16(v)
	}
	return r.ReadUint16()
}

func (r *Reader) ReadInt16Variable() int16 {
	v := r.ReadByte()
	if v < 255 {
		return int16(v) - 127
	}
	return int16(r.ReadUint16())
}

func (r *Reader) ReadUint24() uint32 {
	for r.n < 3 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint32(r.buf[r.at]) | uint32(r.buf[r.at+1])<<8 | uint32(r.buf[r.at+2])<<16
	r.at += 3
	r.n -= 3
	return res
}

func (r *Reader) ReadUint32() uint32 {
	for r.n < 4 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint32(r.buf[r.at]) | uint32(r.buf[r.at+1])<<8 | uint32(r.buf[r.at+2])<<16 | uint32(r.buf[r.at+3])<<24
	r.at += 4
	r.n -= 4
	return res
}

func (r *Reader) ReadUint48() uint64 {
	for r.n < 6 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
	r.at += 6
	r.n -= 6
	return res
}

func (r *Reader) ReadUint64() uint64 {
	for r.n < 8 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	r.at += 8
	r.n -= 8
	return res
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (r *Reader) ReadUint64Variable() uint64 {
	s1 := int(r.ReadByte())
	for r.n < s1 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	var res1 uint64
	switch s1 {
		case 1: res1 = uint64(r.buf[r.at])
		case 2: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += s1
	r.n -= s1
	return res1
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (r *Reader) Read2Uint64sVariable() (uint64, uint64) {
	s2 := r.ReadByte()
	s1 := s2 >> 4
	s2 &= 15
	for r.n < int(s1 + s2) {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	var res1, res2 uint64
	switch s1 {
		case 1: res1 = uint64(r.buf[r.at])
		case 2: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res1 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += int(s1)
	r.n -= int(s1)
	switch s2 {
		case 1: res2 = uint64(r.buf[r.at])
		case 2: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8
		case 3: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16
		case 4: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24
		case 5: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32
		case 6: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40
		case 7: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48
		case 8: res2 = uint64(r.buf[r.at]) | uint64(r.buf[r.at+1])<<8 | uint64(r.buf[r.at+2])<<16 | uint64(r.buf[r.at+3])<<24 | uint64(r.buf[r.at+4])<<32 | uint64(r.buf[r.at+5])<<40 | uint64(r.buf[r.at+6])<<48 | uint64(r.buf[r.at+7])<<56
	}
	r.at += int(s2)
	r.n -= int(s2)
	return res1, res2
}

func (r *Reader) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

func (r *Reader) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}

func (r *Reader) ReadString8() string {
	return string(r.Readx(int(r.ReadByte())))
}

func (r *Reader) ReadString16() string {
	return string(r.Readx(int(r.ReadUint16())))
}

func (r *Reader) ReadString32() string {
	return string(r.Readx(int(r.ReadUint32())))
}

/*
// 12 bits for uint16 and 4 bits for uint8
func (r *Reader) Read12And4() (uint16, uint16) {
	for r.n < 2 {
		copy(r.buf, r.buf[r.at:r.at+r.n])
		r.at = 0
		m, err := r.f.Read(r.buf[r.n:])
		if err != nil {
			panic(err)
		}
		r.n += m
	}
	res := uint16(r.buf[r.at]) | uint16(r.buf[r.at+1])<<8
	r.at += 2
	r.n -= 2
	return res & 4095, res >> 12
}
*/

func (r *Reader) ReadSpecial() (uint8, bool, bool, bool, bool) {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	c := r.buf[r.at]
	var b1, b2, b3, b4 bool
	if c & 128 > 0 {
		b1 = true
	}
	if c & 64 > 0 {
		b2 = true
	}
	if c & 32 > 0 {
		b3 = true
	}
	if c & 16 > 0 {
		b4 = true
	}
	r.at++
	r.n--
	return c & 7, b1, b2, b3, b4
}

func (r *Reader) ReadSpecial2() (uint8, uint8, uint8, bool) {
	for r.n == 0 {
		r.at = 0
		m, err := r.f.Read(r.buf)
		if err != nil {
			panic(err)
		}
		r.n = m
	}
	c := r.buf[r.at]
	var b1 bool
	if c & 128 > 0 {
		b1 = true
	}
	v1 := c & 7
	v2 := (c >> 3) & 3
	v3 := (c >> 5) & 3
	r.at++
	r.n--
	return v1, v2, v3, b1
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if sw, ok := r.f.(io.Seeker); ok {
		r.at, r.n = 0, 0
		return sw.Seek(offset, whence)
	}
	return 0, errors.New(`Does not implement io.Seeker`)
}

func (r *Reader) EOF() error {
	_, err := r.f.Read(r.buf)
	if err == io.EOF {
		return nil
	}
	if err == nil {
		return errors.New(`Not EOF`)
	}
	return err
}

func (r *Reader) Recycle(f io.Reader) {
	r.at, r.n = 0, 0
	r.f = f
}

func (r *Reader) Close() error {
	if sw, ok := r.f.(io.Closer); ok { // Attempt to close underlying writer if it has a Close() method
		return sw.Close()
	}
	return nil
}

// -------- BYTES READER --------

type BytesReader struct {
	data []byte
	cursor, length int
}

func NewBytesReader(p []byte) *BytesReader {
	return &BytesReader{data: p, length: len(p)}
}

func (b *BytesReader) Read(p []byte) (int, error) {
	if b.cursor >= b.length {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.cursor:])
	b.cursor += n
	return n, nil
}

// Readx is like Read but returns a slice of the original instead of copying the bytes into the supplied slice
func (r *BytesReader) Readx(x int) []byte {
	r.cursor += x
	return r.data[r.cursor-x:r.cursor]
}

func (r *BytesReader) ReadBool() bool {
	var b1 bool
	if r.data[r.cursor] > 0 {
		b1 = true
	}
	r.cursor++
	return b1
}

func (r *BytesReader) Read2Bools() (bool, bool) {
	var b1, b2 bool
	switch r.data[r.cursor] {
		case 1: b1 = true
		case 2: b2 = true
		case 3: b1, b2 = true, true
	}
	r.cursor++
	return b1, b2
}

func (r *BytesReader) Read2Unit4s() (uint8, uint8) {
	res1, res2 := r.data[r.cursor] & 15, r.data[r.cursor] >> 4
	r.cursor++
	return res1, res2
}

func (r *BytesReader) ReadByte() uint8 {
	r.cursor++
	return r.data[r.cursor-1]
}

func (r *BytesReader) ReadUTF8() []byte {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return unleak.Bytes(r.data[r.cursor-1:r.cursor])
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		return unleak.Bytes(r.data[r.cursor-2:r.cursor])
	} else {
		r.cursor += 3
		return unleak.Bytes(r.data[r.cursor-3:r.cursor])
	}
}

func (r *BytesReader) ReadUTF8Slice() []byte {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return r.data[r.cursor-1:r.cursor]
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		return r.data[r.cursor-2:r.cursor]
	} else {
		r.cursor += 3
		return r.data[r.cursor-3:r.cursor]
	}
}

func (r *BytesReader) ReadRune() rune {
	if r.data[r.cursor] < 128 { // length 1
		r.cursor++
		return rune(r.data[r.cursor-1])
	}
	if r.data[r.cursor] & 32 == 0 { // length 2
		r.cursor += 2
		r, _ := utf8.DecodeRune(r.data[r.cursor-2:r.cursor])
		return r
	} else {
		r.cursor += 3
		r, _ := utf8.DecodeRune(r.data[r.cursor-3:r.cursor])
		return r
	}
}

func (r *BytesReader) ReadUint16() uint16 {
	r.cursor += 2
	return uint16(r.data[r.cursor-2]) | uint16(r.data[r.cursor-1])<<8
}

// If it's less than 255 then it's encoded in the 1st byte, otherwise 1st byte is 255 and it's encoded in two more bytes
// This is only useful if it is expected that the value will be <255 more than half the time
func (r *BytesReader) ReadUint16Variable() uint16 {
	v := r.ReadByte()
	if v < 255 {
		return uint16(v)
	}
	return r.ReadUint16()
}

func (r *BytesReader) ReadInt16Variable() int16 {
	v := r.ReadByte()
	if v < 255 {
		return int16(v) - 127
	}
	return int16(r.ReadUint16())
}

func (r *BytesReader) ReadUint24() uint32 {
	r.cursor += 3
	return uint32(r.data[r.cursor-3]) | uint32(r.data[r.cursor-2])<<8 | uint32(r.data[r.cursor-1])<<16
}

func (r *BytesReader) ReadUint32() uint32 {
	r.cursor += 4
	return uint32(r.data[r.cursor-4]) | uint32(r.data[r.cursor-3])<<8 | uint32(r.data[r.cursor-2])<<16 | uint32(r.data[r.cursor-1])<<24
}

func (r *BytesReader) ReadUint48() uint64 {
	r.cursor += 6
	return uint64(r.data[r.cursor-6]) | uint64(r.data[r.cursor-5])<<8 | uint64(r.data[r.cursor-4])<<16 | uint64(r.data[r.cursor-3])<<24 | uint64(r.data[r.cursor-2])<<32 | uint64(r.data[r.cursor-1])<<40
}

func (r *BytesReader) ReadUint64() uint64 {
	r.cursor += 8
	return uint64(r.data[r.cursor-8]) | uint64(r.data[r.cursor-7])<<8 | uint64(r.data[r.cursor-6])<<16 | uint64(r.data[r.cursor-5])<<24 | uint64(r.data[r.cursor-4])<<32 | uint64(r.data[r.cursor-3])<<40 | uint64(r.data[r.cursor-2])<<48 | uint64(r.data[r.cursor-1])<<56
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (r *BytesReader) ReadUint64Variable() uint64 {
	s1 := int(r.ReadByte())
	var res1 uint64
	switch s1 {
		case 1: res1 = uint64(r.data[r.cursor])
		case 2: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += s1
	return res1
}

// The first byte stores the bit length of the two integers. Then come the two integers. Length is only 1 byte more than the smallest representation of both integers.
func (r *BytesReader) Read2Uint64sVariable() (uint64, uint64) {
	s2 := r.ReadByte()
	s1 := s2 >> 4
	s2 &= 15
	var res1, res2 uint64
	switch s1 {
		case 1: res1 = uint64(r.data[r.cursor])
		case 2: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res1 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += int(s1)
	switch s2 {
		case 1: res2 = uint64(r.data[r.cursor])
		case 2: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8
		case 3: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16
		case 4: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24
		case 5: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32
		case 6: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40
		case 7: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48
		case 8: res2 = uint64(r.data[r.cursor]) | uint64(r.data[r.cursor+1])<<8 | uint64(r.data[r.cursor+2])<<16 | uint64(r.data[r.cursor+3])<<24 | uint64(r.data[r.cursor+4])<<32 | uint64(r.data[r.cursor+5])<<40 | uint64(r.data[r.cursor+6])<<48 | uint64(r.data[r.cursor+7])<<56
	}
	r.cursor += int(s2)
	return res1, res2
}

func (r *BytesReader) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

func (r *BytesReader) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}

func (r *BytesReader) ReadString8() string {
	return string(r.Readx(int(r.ReadByte())))
}

func (r *BytesReader) ReadString16() string {
	return string(r.Readx(int(r.ReadUint16())))
}

func (r *BytesReader) ReadString32() string {
	return string(r.Readx(int(r.ReadUint32())))
}

/*
// 12 bits for uint16 and 4 bits for uint8
func (r *BytesReader) Read12And4() (uint16, uint16) {
	res := uint16(r.data[r.cursor]) | uint16(r.data[r.cursor+1])<<8
	r.cursor += 2
	return res & 4095, res >> 12
}
*/

func (r *BytesReader) ReadSpecial() (uint8, bool, bool, bool, bool) {
	c := r.data[r.cursor]
	var b1, b2, b3, b4 bool
	if c & 128 > 0 {
		b1 = true
	}
	if c & 64 > 0 {
		b2 = true
	}
	if c & 32 > 0 {
		b3 = true
	}
	if c & 16 > 0 {
		b4 = true
	}
	r.cursor++
	return c & 7, b1, b2, b3, b4
}

func (r *BytesReader) ReadSpecial2() (uint8, uint8, uint8, bool) {
	c := r.data[r.cursor]
	var b1 bool
	if c & 128 > 0 {
		b1 = true
	}
	v1 := c & 7
	v2 := (c >> 3) & 3
	v3 := (c >> 5) & 3
	r.cursor++
	return v1, v2, v3, b1
}

func (b *BytesReader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
		case 0:
			abs = offset
		case 1:
			abs = int64(b.cursor) + offset
		case 2:
			abs = int64(b.length) + offset
		default:
			return 0, errors.New("buffer.BytesReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("buffer.BytesReader.Seek: negative position")
	}
	b.cursor = int(abs)
	return abs, nil
}

func (r *BytesReader) EOF() error {
	if r.cursor == len(r.data) {
		return nil
	}
	return errors.New(`Not EOF`)
}
