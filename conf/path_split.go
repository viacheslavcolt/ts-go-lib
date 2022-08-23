package conf

const sep = '.'

type _PathSpliter struct {
	last int
	len  int

	path []byte
}

func (ps *_PathSpliter) _Init(path string) {
	ps.last = 0
	ps.path = []byte(path)

	ps.len = len(ps.path)
}

func (ps *_PathSpliter) _Next() []byte {
	var (
		i     int
		slice []byte
	)

	i = ps.last

	for ; ps.last < ps.len; ps.last++ {
		if ps.path[ps.last] == sep {
			slice = ps.path[i:ps.last]
			ps.last++
			return slice
		}
	}

	if i < ps.len {
		return ps.path[i:ps.last]
	}

	return nil
}

func (ps *_PathSpliter) _CanTakeNext() bool {
	return ps.last < ps.len
}
