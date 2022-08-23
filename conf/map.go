package conf

/*
 */

type _EntryType int

const (
	_MapPtrT _EntryType = iota
	_ValEntryT
)

type _map map[string]_MapEntry

func (m _map) exists(key string) bool {
	var (
		ok bool
	)

	_, ok = m[key]

	return ok
}

func (m _map) setObj(obj *_TscObj) {
	var (
		entry  _MapEntry
		newMap _map
	)

	for obj != nil {
		if obj.T == _ObjNamespaceT {
			newMap = make(_map)

			entry = _MapEntry{
				T:      _MapPtrT,
				Obj:    nil,
				MapPtr: newMap,
			}

			m[string(obj.Key)] = entry

			newMap.setObj(obj.Child)
		} else {
			entry = _MapEntry{
				T:      _ValEntryT,
				Obj:    obj,
				MapPtr: nil,
			}

			m[string(obj.Key)] = entry
		}

		obj = obj.Next
	}
}

type _MapEntry struct {
	T      _EntryType
	Obj    *_TscObj
	MapPtr _map
}

type _CfgMap struct {
	root _map
}

func _NewCfgMap() *_CfgMap {
	return &_CfgMap{
		root: make(map[string]_MapEntry),
	}
}

func (c *_CfgMap) _SetRoot(root *_TscObj) {
	c.root.setObj(root.Child)
}

func (c *_CfgMap) _Get(path string) *_TscObj {
	var (
		splitter  _PathSpliter
		pathShard []byte
		key       string
		m         _map
		entry     _MapEntry
	)

	m = c.root

	splitter._Init(path)

	pathShard = splitter._Next()

	for pathShard != nil {
		key = string(pathShard)

		if !m.exists(key) {
			return nil
		}

		entry = m[key]

		if splitter._CanTakeNext() {
			if entry.T != _MapPtrT {
				return nil
			}

			m = entry.MapPtr
		} else {
			if entry.T != _ValEntryT {
				return nil
			}

			return entry.Obj
		}

		pathShard = splitter._Next()
	}

	return nil
}
