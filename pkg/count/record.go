package count

import (
	"sort"
	"time"
)

var TimeLedger map[string]time.Time = make(map[string]time.Time, 0)

type Record struct {
	File string
	Line int
	Err  error
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
		l[r.File] = l[r.File].Incriment()
	} else {
		l[r.File] = Count{
			Record: r,
			Count:  1,
		}
	}

	if _, ok := TimeLedger[r.File]; !ok {
		TimeLedger[r.File] = time.Now()
	}
}

func (l Ledger) AddCount(f string, c Count) {
	if _, ok := l[f]; ok {
		l[f] = l[f].AddCount(c)
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
	switch {
	case len(l) == 0:
		return []string{}
	case len(l) < c:
		c = len(l)
	}

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

func (l Ledgers) GetPrev() Ledger {
	if len(l) < 2 {
		return NewLedger()
	}

	return l[len(l)-2]
}

func (l Ledgers) GetLast() Ledger {
	if len(l) < 1 {
		return NewLedger()
	}

	return l[len(l)-1]
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

func (c Count) AddCount(nC Count) Count {
	c.Count = c.Count + nC.Count

	return c
}

func (c Count) Incriment() Count {
	c.Count++

	return c
}
