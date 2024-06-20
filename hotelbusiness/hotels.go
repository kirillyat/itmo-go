//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {

	d := make(map[int]int)

	for _, g := range guests {
		d[g.CheckInDate]++
		d[g.CheckOutDate]--
	}

	keys := make([]int, 0, len(d))

	for k := range d {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var cnt int
	ans := make([]Load, 0, len(d))

	for _, k := range keys {
		if d[k] == 0 {
			continue
		}
		cnt += d[k]
		ans = append(ans, Load{StartDate: k, GuestCount: cnt})
	}

	return ans
}
