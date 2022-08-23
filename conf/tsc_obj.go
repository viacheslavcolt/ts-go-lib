package conf

type _ObjT uint

const (
	_ObjStringT _ObjT = iota
	_ObjIntT
	_ObjBoolT
	_ObjNamespaceT
)

type _TscObj struct {
	T _ObjT

	Root  *_TscObj
	Next  *_TscObj
	Child *_TscObj

	Key    []byte
	StrVal []byte
	Int    int
	Bool   bool
}

func _NewTscObj() *_TscObj {
	return &_TscObj{
		T:      0,
		Root:   nil,
		Next:   nil,
		Child:  nil,
		Key:    nil,
		StrVal: nil,
		Int:    0,
		Bool:   false,
	}
}
