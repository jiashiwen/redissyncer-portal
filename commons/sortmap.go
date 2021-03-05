package commons

import "sort"

type Pair struct {
	Key   string
	Value int64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// @title    SortMapByValue
// @description   map 按 value 排序
// @auth      JiaShiwen      2021/3/   10:57
// @param     sortedMap      map[string]int64        "需要排序的map"
// @param     reverse        reverse       "反向排序 "
// @return    PairList        PairList         "键值对数组"
func SortMapByValue(sortedMap map[string]int64, reverse bool) PairList {
	pl := make(PairList, len(sortedMap))
	i := 0
	for k, v := range sortedMap {
		pl[i] = Pair{k, v}
		i++
	}
	if reverse {
		sort.Sort(sort.Reverse(pl))
		return pl
	}
	sort.Sort(pl)
	return pl
}
