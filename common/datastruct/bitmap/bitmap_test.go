package bitmap

import (
	"testing"
)

func TestNew(t *testing.T) {
	bitmap := NewBitMap()
	if bitmap == nil || len(*bitmap) != 0 {
		t.Errorf("New() failed, expected empty BitMap")
	}
}

func TestToByteSize(t *testing.T) {
	if toByteSize(8) != 1 {
		t.Errorf("toByteSize() failed, expected 1 for 8 bits")
	}
	if toByteSize(9) != 2 {
		t.Errorf("toByteSize() failed, expected 2 for 9 bits")
	}
}

func TestBitMap_Grow(t *testing.T) {
	bitmap := NewBitMap()
	bitmap.grow(16)
	if len(*bitmap) != 2 {
		t.Errorf("grow() failed, expected length 2, got %d", len(*bitmap))
	}
}

func TestBitMap_BitSize(t *testing.T) {
	bitmap := NewBitMap()
	bitmap.grow(16)
	if bitmap.BitSize() != 16 {
		t.Errorf("BitSize() failed, expected 16, got %d", bitmap.BitSize())
	}
}

func TestFromBytes(t *testing.T) {
	bytes := []byte{0xFF, 0x00}
	bitmap := FromBytes(bytes)
	if len(*bitmap) != 2 || (*bitmap)[0] != 0xFF || (*bitmap)[1] != 0x00 {
		t.Errorf("FromBytes() failed, expected [0xFF, 0x00]")
	}
}

func TestBitMap_ToBytes(t *testing.T) {
	bitmap := FromBytes([]byte{0xFF, 0x00})
	bytes := bitmap.ToBytes()
	if len(bytes) != 2 || bytes[0] != 0xFF || bytes[1] != 0x00 {
		t.Errorf("ToBytes() failed, expected [0xFF, 0x00]")
	}
}

func TestBitMap_SetBit(t *testing.T) {
	bitmap := NewBitMap()
	bitmap.SetBit(10, 1)
	if bitmap.GetBit(10) != 1 {
		t.Errorf("SetBit() failed, expected bit 10 to be 1")
	}
	bitmap.SetBit(10, 0)
	if bitmap.GetBit(10) != 0 {
		t.Errorf("SetBit() failed, expected bit 10 to be 0")
	}
}

func TestBitMap_GetBit(t *testing.T) {
	bitmap := NewBitMap()
	if bitmap.GetBit(5) != 0 {
		t.Errorf("GetBit() failed, expected bit 5 to be 0")
	}
	bitmap.SetBit(5, 1)
	if bitmap.GetBit(5) != 1 {
		t.Errorf("GetBit() failed, expected bit 5 to be 1")
	}
}

func TestBitMap_ForEachBit(t *testing.T) {
	bitmap := NewBitMap()
	bitmap.SetBit(0, 1)
	bitmap.SetBit(1, 0)
	bitmap.SetBit(2, 1)

	count := 0
	bitmap.ForEachBit(0, 3, func(offset int64, val byte) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("ForEachBit() failed, expected 3 iterations, got %d", count)
	}
}

func TestBitMap_ForEachByte(t *testing.T) {
	bitmap := FromBytes([]byte{0xFF, 0x00})
	count := 0
	bitmap.ForEachByte(0, 2, func(offset int64, val byte) bool {
		count++
		return true
	})
	if count != 2 {
		t.Errorf("ForEachByte() failed, expected 2 iterations, got %d", count)
	}
}
