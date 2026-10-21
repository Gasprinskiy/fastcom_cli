package main

type FixedByteBuf struct {
	Index int8
	BSize int8
	ASize int8
	Arr   []float64
}

func NewFixedByteBuf(size int8) *FixedByteBuf {
	return &FixedByteBuf{
		Index: -1,
		BSize: size,
		Arr:   make([]float64, size),
	}
}

func (b *FixedByteBuf) Add(val float64) {
	b.Index += 1
	if b.Index >= b.BSize {
		b.Index = 0
	}
	if b.ASize < b.BSize {
		b.ASize++
	}

	b.Arr[b.Index] = val
}
