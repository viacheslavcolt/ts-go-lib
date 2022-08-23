package conf

import (
	"os"
)

type _Buffer struct {
	src    []byte
	size   uint
	offset uint
}

func _InitBuffer(buf *_Buffer, filePath string) error {
	var (
		err error
	)

	if buf.src, err = os.ReadFile(filePath); err != nil {
		return err
	}

	buf.offset = 0
	buf.size = uint(len(buf.src))

	return nil
}

func (b *_Buffer) increasePosition(val uint) {
	b.offset += val
}

func (b *_Buffer) canAccessAtOffset(index uint) bool {
	if (b.offset + index) < b.size {
		return true
	}

	return false
}

func (b *_Buffer) _CanRead(size uint) bool {
	return b.canAccessAtOffset(size - 1)
}

func (b *_Buffer) _AdvanceAt(size uint) {
	b.increasePosition(size)
}

func (b *_Buffer) _CanAccessAtIndex(index uint) bool {
	return b.canAccessAtOffset(index)
}

func (b *_Buffer) _AtIndex(index uint) byte {
	return b.src[b.offset+index]
}

func (b *_Buffer) _CmpCharAtIndex(index uint, char byte) bool {
	return b._AtIndex(index) == char
}

func (b *_Buffer) _GetOffset() uint {
	return b.offset
}

func (b *_Buffer) _CmpStringAtIndex(index uint, str []byte) bool {
	var (
		i uint
		l uint
	)

	i = 0
	l = uint(len(str))

	for i = 0; i < l; i++ {
		if b._AtIndex(index+i) != str[i] {
			return false
		}
	}

	return true
}

// skip all spaces, all characters up to code 32 in ascii table
func (b *_Buffer) _SkipWhiteSpace() bool {
	for b._CanAccessAtIndex(0) {
		if isASpace(b._AtIndex(0)) {
			b._AdvanceAt(1)
			continue
		}

		break
	}

	return true
}

func (b *_Buffer) _MakeSlice(start uint, end uint) []byte {
	return b.src[b.offset+start : b.offset+end]
}

func isASpace(char byte) bool {
	return char <= 32
}
