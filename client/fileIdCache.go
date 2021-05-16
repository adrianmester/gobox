package main

/*
pathIDMap allocates a unique numeric ID to each path
 */
type pathIDMap struct {
	mp        map[string]int64
	highestID int64
}

func NewPathIDMap() *pathIDMap {
	return &pathIDMap{
		map[string]int64{},
		0,
	}
}

/*
GetID returns the numeric ID associated with the give path. If it doesn't exist
a new ID is created. Uniqueness is guaranteed by always increasing the `highestID` member
of the struct.
 */
func (p *pathIDMap) GetID(path string) int64 {
	fid, ok := p.mp[path]
	if ok {
		return fid
	}
	p.highestID += 1
	p.mp[path] = p.highestID
	return p.highestID
}
