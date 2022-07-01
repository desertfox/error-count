package count

import (
	"sort"
	"time"
)

type Record struct {
	File    string
	Line    int
	Created time.Time
	Err     error
}

type Count struct {
	Record Record
	Count  int
}

type Ledger map[string]Count

type Ledgers []Ledger

func NewLedger() Ledger {
	return make(Ledger, 0)
}

func (l Ledger) Add(r Record) {
	if _, ok := l[r.File]; ok {
		count := l[r.File]
		count.Count++
		l[r.File] = count
	} else {
		l[r.File] = Count{
			Record: r,
			Count:  0,
		}
	}
}

func (l Ledger) AddCount(f string, c Count) {
	if v, ok := l[f]; ok {
		v.Count = v.Count + c.Count
	} else {
		l[f] = c
	}
}

func (l Ledger) GetCount(file string) Count {
	if c, ok := l[file]; ok {
		return c
	}
	return Count{
		Record: Record{
			File: file,
		},
		Count: 0,
	}
}

func (l Ledger) GetTopFileInstances(c int) []string {
	keys := make([]string, 0, len(l))
	for key := range l {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return l[keys[i]].Count > l[keys[j]].Count })

	return keys[0:c]
}

func (l *Ledgers) Add(newL Ledger) {
	*l = append(*l, newL)
}

func (l Ledgers) GetLast() Ledger {
	if len(l) == 0 {
		return NewLedger()
	}

	return l[0]
}

func (l Ledgers) TotalLedger() Ledger {
	nl := NewLedger()

	for _, ledger := range l {
		for file, count := range ledger {
			nl.AddCount(file, count)
		}
	}

	return nl
}
