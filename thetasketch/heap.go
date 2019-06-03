package thetasketch

type UintHeap struct {
	data     []uint64
	dict     map[uint64]struct{}
	bound    int
}

func NewHeap(n int) UintHeap {
	return UintHeap{make([]uint64, 0), make(map[uint64]struct{}), n}
}

func (h *UintHeap) Peak() uint64 {
	if h.Len() != 0 {
		return h.data[0]
	} else {
		return 0
	}
}

func (h *UintHeap) shiftUp() {
	i := h.Len()-1
	for i>0 {
		if h.data[i] < h.data[(i-1)/2] {
			h.data[i], h.data[(i-1)/2] = h.data[(i-1)/2], h.data[i]
			i = (i-1)/2
		} else if h.data[i] > h.data[(i-1)/2] {
			return
		} else {
			return
		}
	}
}

func (h *UintHeap) shiftDown() {
	i := 0
	for i<(h.Len()-1)/2 {
		if h.data[i] < h.data[i*2+1] && h.data[i] < h.data[i*2+2] {
			return
		} else if h.data[i] > h.data[i*2+1] && h.data[i*2+2] > h.data[i*2+1] {
			h.data[i], h.data[i*2+1] = h.data[i*2+1], h.data[i]
			i = i*2+1
		} else if h.data[i*2+2] < h.data[i] && h.data[i*2+2] < h.data[i*2+1] {
			h.data[i], h.data[i*2+2] = h.data[i*2+2], h.data[i]
			i = i*2+2
		} else {
			return
		}
	}
}

func (h *UintHeap) insertOne(n uint64) {
	if _, ok := h.dict[n]; ok {
		return
	}
	if n <= h.Peak() && h.Full() {
		return
	}
	if h.Len() == h.bound {
		h.Pop()
	}
	h.data = append(h.data, n)
	h.dict[n] = struct{}{}
	h.shiftUp()
}

func (h *UintHeap) Push(n uint64) {
	h.insertOne(n)
}

func (h *UintHeap) Pop() uint64 {
	if h.Len() == 0 {
		return 0
	}
	x := h.data[0]
	delete(h.dict, x)
	if h.Len() == 1 {
		h.data = h.data[1:]
	} else {
		h.data[0] = h.data[h.Len()-1]
		h.data = h.data[:h.Len()-1]
		h.shiftDown()
	}
	return x
}

func (h *UintHeap) Len() int {
	return len(h.data)
}

func (h *UintHeap) Full() bool {
	return h.Len() == h.bound
}

func (h *UintHeap) Copy() UintHeap {
	dst := make([]uint64, len(h.data))
	dict := make(map[uint64]struct{})
	copy(dst, h.data)
	for k, _ := range h.dict {
		dict[k] = struct{}{}
	}
	return UintHeap{dst, dict, h.bound}
}

func (h *UintHeap) Items() []uint64 {
	dst := make([]uint64, len(h.data))
	copy(dst, h.data)
	return dst
}