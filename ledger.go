package main

import (
	"sort"
	"time"
)

type Record struct {
	File string
	Line int
	Err  error
}

type Count struct {
	Record Record
	Count  int
}

type TimeLedger map[string]time.Time

type Ledger map[string]Count

type Ledgers []Ledger

func (tL TimeLedger) Add(r Record) {
	if _, ok := tL[r.File]; !ok {
		tL[r.File] = time.Now()
	}
}

func (l Ledger) Incriment(r Record) {
	if _, ok := l[r.File]; ok {
		l[r.File] = l[r.File].Incriment()
	} else {
		l[r.File] = Count{
			Record: r,
			Count:  1,
		}
	}
}

func (l Ledger) Add(f string, c Count) {
	if _, ok := l[f]; ok {
		l[f] = l[f].Add(c)
	} else {
		l[f] = c
	}
}

func (l Ledger) Get(file string) Count {
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

func (l Ledger) Top(c int) []string {
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

func (l Ledgers) Prev() Ledger {
	if len(l) < 2 {
		return make(Ledger, 0)
	}

	return l[len(l)-2]
}

func (l Ledgers) Last() Ledger {
	if len(l) < 1 {
		return make(Ledger, 0)
	}

	return l[len(l)-1]
}

func (l Ledgers) Total() Ledger {
	newLedger := make(Ledger, 0)

	for _, ledger := range l {
		for file, count := range ledger {
			newLedger.Add(file, count)
		}
	}

	return newLedger
}

func (c Count) Add(nC Count) Count {
	c.Count = c.Count + nC.Count

	return c
}

func (c Count) Incriment() Count {
	c.Count++

	return c
}
