package rcache

import (
	"fmt"
	"sync"
)

type Group struct {
	Name    string
	FleetID int
}

type Groups struct {
	Data map[string]Group
	sync.RWMutex
}

var (
	Grouplist Groups
)

func (g *Groups) Put(id string, group Group) {
	if len(g.Data) == 0 {
		g.Data = make(map[string]Group)
	}
	g.Lock()
	g.Data[id] = group
	g.Unlock()
}

func (g *Groups) Get(id string) (Group, error) {
	g.Lock()
	group, ok := g.Data[id]
	if !ok {
		g.Unlock()
		return group, fmt.Errorf("group (id %s) not found", id)
	}
	g.Unlock()
	return group, nil
}
