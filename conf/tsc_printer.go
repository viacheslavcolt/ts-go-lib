package conf

import (
	"fmt"
	"strconv"
)

func _TscPrint(obj *_TscObj, deep int) {
	if obj.Key == nil {
		fmt.Println("Namespace: root")
	} else {
		fmt.Printf("%sNamespace: %s\n", getIndent(deep-1), string(obj.Key))
	}

	var (
		val string
	)

	for o := obj.Child; o != nil; o = o.Next {
		switch o.T {
		case _ObjNamespaceT:
			print(o, deep+1)
			continue
		case _ObjBoolT:
			val = boolToString(o.Bool)
			break
		case _ObjStringT:
			val = string(o.StrVal)
			break
		case _ObjIntT:
			val = strconv.Itoa(o.Int)
			break
		}
		fmt.Printf("%skey: %s, type: %s, val: %s\n", getIndent(deep), string(o.Key), getType(o.T), val)
	}
}

func getIndent(deep int) string {
	if deep <= 0 {
		return ""
	}

	indent := make([]byte, 0, deep)

	for deep > 0 {
		indent = append(indent, '\t')
		deep--
	}

	return string(indent)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}

	return "false"
}

func getType(t _ObjT) string {
	switch t {
	case _ObjIntT:
		return "int"
	case _ObjStringT:
		return "string"
	case _ObjBoolT:
		return "bool"
	}
	return "unknow"
}
