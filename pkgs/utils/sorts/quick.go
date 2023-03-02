package sorts

type Item interface {
	Greater(Item) bool
	String() string
}

func QSort(iList []Item) (r []Item) {
	if len(iList) < 1 {
		return iList
	}
	mid := iList[0]
	left := make([]Item, 0)
	right := make([]Item, 0)
	for i := 1; i < len(iList); i++ {
		if mid.Greater(iList[i]) {
			left = append(left, iList[i])
		} else {
			right = append(right, iList[i])
		}
	}
	left, right = QSort(left), QSort(right)
	r = append(r, left...)
	r = append(r, mid)
	r = append(r, right...)
	return r
}

func QuickSort(iList []Item) (r []string) {
	items := QSort(iList)
	for _, itm := range items {
		r = append(r, itm.String())
	}
	return
}
