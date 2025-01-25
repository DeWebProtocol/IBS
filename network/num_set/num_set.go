package num_set

type Set struct {
	s []uint64
}

func NewSet() *Set {
	return &Set{
		[]uint64{},
	}
}

func (t *Set) Insert(num uint64) bool {
	l := t.locate(num)
	if len(t.s) > 0 && t.s[l] == num {
		return false
	}
	t.s = append(t.s, num)
	for i := len(t.s) - 1; i > 0; i-- {
		if t.s[i-1] > t.s[i] {
			t.s[i], t.s[i-1] = t.s[i-1], t.s[i]
		}
	}
	return true
}

// TODO: use bin-search
func (t *Set) locate(num uint64) int {
	for i := len(t.s) - 1; i >= 0; i-- {
		if t.s[i] <= num {
			return i
		}
	}
	return 0
	//left, right := 0, len(t.s)-1
	//for left < right {
	//	mid := (left + right) / 2
	//	if num > t.s[mid] {
	//		left = mid + 1
	//	} else {
	//		right = mid - 1
	//	}
	//}
	//if left < len(t.s) {
	//	return left
	//}
	//return left - 1
}

// Around : n numbers around num
func (t *Set) Around(num uint64, n int) []uint64 {
	pos := t.locate(num)
	//fmt.Println("locate", num, "at", pos)
	left := pos - n/2
	right := pos + n/2 + 1
	if left < 0 {
		left = 0
	}
	if right > len(t.s) {
		right = len(t.s)
	}
	return t.s[left:right]
}
