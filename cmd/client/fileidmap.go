package main

import "sync"

/*
pathIDMap allocates a unique numeric ID to each path
 */
type pathIDMap struct {
	mp        map[string]int64
	highestID int64
	mu		  *sync.Mutex
}

func newPathIDMap() *pathIDMap {
	return &pathIDMap{
		map[string]int64{},
		0,
		&sync.Mutex{},
	}
}

/*
GetID returns the numeric ID associated with the give path. If it doesn't exist
a new ID is created. Uniqueness is guaranteed by always increasing the `highestID` member
of the struct.
 */
func (p *pathIDMap) GetID(path string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	fid, ok := p.mp[path]
	if ok {
		return fid
	}
	p.highestID += 1
	p.mp[path] = p.highestID
	return p.highestID
}
