package tsc

type ObjT uint

const (
	ObjStringT ObjT = iota
	ObjIntT
	ObjBoolT
	ObjNamespaceT
)

type TscObj struct {
	T ObjT

	Root  *TscObj
	Next  *TscObj
	Child *TscObj

	Key    []byte
	StrVal []byte
	Int    int
	Bool   bool
}

func _NewTscObj() *TscObj {
	return &TscObj{
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
