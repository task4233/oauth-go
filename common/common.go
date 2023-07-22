package common

func AreTwoUnorderedSlicesSame[E comparable](s []E, t []E) bool {
	if len(s) != len(t) {
		return false
	}

	mpS := map[E]int{}
	for _, e := range s {
		mpS[e]++
	}
	for _, e := range t {
		mpS[e]--
	}
	for _, v := range mpS {
		if v != 0 {
			return false
		}
	}

	return true
}
