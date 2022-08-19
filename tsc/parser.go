package tsc

import (
	"fmt"
)

func Parse(filePath string) (*TscObj, error) {
	var (
		root *TscObj
		buf  _Buffer
		err  error
	)

	_InitBuffer(&buf)

	if err = buf._Load(filePath); err != nil {
		return nil, err
	}

	root = _NewTscObj()

	if err = parseNamespace(root, &buf, true); err != nil {
		return nil, err
	}

	return root, nil
}

func parseNamespace(obj *TscObj, buf *_Buffer, isRoot bool) error {
	var (
		newObj        *TscObj
		nextPtr       **TscObj
		isReachingEnd bool
		err           error
	)

	isReachingEnd = false

	if !isRoot {
		if !buf._CanAccessAtIndex(0) || !buf._CmpCharAtIndex(0, '{') {
			return fmt.Errorf("expected { on the position %d", buf._GetOffset())
		}

		buf._AdvanceAt(1)
	}

	obj.T = ObjNamespaceT

	nextPtr = &obj.Child

	for buf._SkipWhiteSpace() && buf._CanAccessAtIndex(0) {
		if !isRoot {
			if buf._CmpCharAtIndex(0, '}') {
				buf._AdvanceAt(1)
				isReachingEnd = true
				break
			}
		}

		if skipComment(buf) {
			continue
		}

		newObj = _NewTscObj()

		if err = parsekey(newObj, buf); err != nil {
			return err
		}

		//fmt.Println(newObj)
		buf._SkipWhiteSpace()

		if !buf._CanAccessAtIndex(0) {
			return fmt.Errorf("unexpected end of file on position %d", buf._GetOffset())
		}

		switch buf._AtIndex(0) {
		case '{':
			if err = parseNamespace(newObj, buf, false); err != nil {
				return err
			}
			break
		case '=':
			buf._AdvanceAt(1)

			buf._SkipWhiteSpace()

			if err = parseVal(newObj, buf); err != nil {
				return err
			}

			if !buf._CanAccessAtIndex(0) {
				return fmt.Errorf("unexpected end of file on position %d", buf._GetOffset())
			}

			if buf._AtIndex(0) != ';' {
				return fmt.Errorf("expected character ; on the postition %d", buf._GetOffset())
			}

			buf._AdvanceAt(1)

			break
		default:
			return fmt.Errorf("unexpected char %c on the position %d", buf._AtIndex(0), buf._GetOffset())
		}

		// set pointers for valid struct of parsed data
		newObj.Root = obj
		*nextPtr = newObj
		nextPtr = &(*nextPtr).Next
	}

	if !isRoot && !isReachingEnd {
		return fmt.Errorf("expected character } on the postition %d", buf._GetOffset())
	}

	return nil
}

/*

{
	=f
}

*/
func parsekey(obj *TscObj, buf *_Buffer) error {
	var (
		startIndex    int
		endIndex      int
		isInterrapted bool

		allocSlice []byte
	)

	isInterrapted = false
	startIndex = 0
	endIndex = startIndex

	for ; buf._CanAccessAtIndex(uint(endIndex)); endIndex++ {
		if !_isAvailableForKey(buf._AtIndex(uint(endIndex))) {
			isInterrapted = true
			break
		}
	}

	if isInterrapted {
		if startIndex == endIndex {
			return fmt.Errorf("unexpected %c on the position %d, expected start of key name", buf._AtIndex(0), buf._GetOffset())
		}

		allocSlice = make([]byte, (endIndex - startIndex))

		copy(allocSlice, buf._MakeSlice(uint(startIndex), uint(endIndex+1)))

		buf._AdvanceAt(uint(endIndex))

		obj.Key = allocSlice

		return nil
	}

	return fmt.Errorf("unexpected end")
}

func parseVal(obj *TscObj, buf *_Buffer) error {
	if !buf._CanAccessAtIndex(0) {
		return fmt.Errorf("expected start value of key %s on the position %d", string(obj.Key), buf._GetOffset())
	}

	if buf._CanRead(4) && buf._CmpStringAtIndex(0, []byte("true")) {
		buf._AdvanceAt(4)

		obj.T = ObjBoolT
		obj.Bool = true

		return nil
	}

	if buf._CanRead(5) && buf._CmpStringAtIndex(0, []byte("false")) {
		buf._AdvanceAt(5)

		obj.T = ObjBoolT
		obj.Bool = false

		return nil
	}

	if buf._CmpCharAtIndex(0, '"') {
		return parseStringVal(obj, buf)
	}

	if buf._CmpCharAtIndex(0, '-') || _isNumber(buf._AtIndex(0)) {

	}

	return nil
}

func parseStringVal(obj *TscObj, buf *_Buffer) error {
	if !buf._CmpCharAtIndex(0, '"') {
		return fmt.Errorf("expected start of string with \" character, on the position %d", buf._GetOffset())
	}

	buf._AdvanceAt(1)

	return nil
	// var (
	// 	startIndex int
	// 	endIndex   int

	// 	allocSize int
	// )

	// startIndex = 0
	// endIndex = startIndex

	// for ; buf._CanAccessAtIndex(uint(endIndex)) && buf._AtIndex(uint(endIndex)) != '"'; endIndex++ {
	// 	if buf._AtIndex(uint(endIndex)) == '\\' {
	// 		endIndex++

	// 	}
	// }

	// if buf._AtIndex(uint(endIndex)) != '"' {
	// 	return fmt.Errorf("unexpected end")
	// }

	// if startIndex == endIndex {

	// }

	// switch {

	// }
}

func skipComment(buf *_Buffer) bool {
	var (
		isReachingBackspace bool
	)

	isReachingBackspace = false

	if buf._CanRead(2) && buf._CmpStringAtIndex(0, []byte("//")) {
		buf._AdvanceAt(2)

		for buf._CanAccessAtIndex(0) && !isReachingBackspace {
			if buf._CmpCharAtIndex(0, '\n') {
				isReachingBackspace = true
			}

			buf._AdvanceAt(1)
		}

		return true
	}

	return false
}

// available only number, uppercase character, lowercase character for key
func _isAvailableForKey(char byte) bool {
	if (char >= 97 && char <= 122) || (char >= 65 && char <= 90) || (char >= 48 && char <= 57) {
		return true
	}

	return false
}

func _isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}
