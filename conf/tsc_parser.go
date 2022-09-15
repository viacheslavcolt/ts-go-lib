package conf

import (
	"fmt"
	"strconv"
)

func _TscParse(buf *_Buffer) (*_TscObj, error) {
	var (
		root *_TscObj
		err  error
	)

	root = _NewTscObj()

	if err = parseNamespace(root, buf, true); err != nil {
		return nil, err
	}

	return root, nil
}

func parseNamespace(obj *_TscObj, buf *_Buffer, isRoot bool) error {
	var (
		newObj        *_TscObj
		nextPtr       **_TscObj
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

	obj.T = _ObjNamespaceT

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

		newObj.Root = obj
		*nextPtr = newObj
		nextPtr = &(*nextPtr).Next
	}

	if !isRoot && !isReachingEnd {
		return fmt.Errorf("expected character } on the postition %d", buf._GetOffset())
	}

	return nil
}

func parsekey(obj *_TscObj, buf *_Buffer) error {
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
		if !isKeyCharacter(buf._AtIndex(uint(endIndex))) {
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

func parseVal(obj *_TscObj, buf *_Buffer) error {
	if !buf._CanAccessAtIndex(0) {
		return fmt.Errorf("expected start value of key %s on the position %d", string(obj.Key), buf._GetOffset())
	}

	if buf._CanRead(4) && buf._CmpStringAtIndex(0, []byte("true")) {
		buf._AdvanceAt(4)

		obj.T = _ObjBoolT
		obj.Bool = true

		return nil
	}

	if buf._CanRead(5) && buf._CmpStringAtIndex(0, []byte("false")) {
		buf._AdvanceAt(5)

		obj.T = _ObjBoolT
		obj.Bool = false

		return nil
	}

	if buf._CmpCharAtIndex(0, '"') {
		return parseStringVal(obj, buf)
	}

	if buf._CmpCharAtIndex(0, '-') || isANumber(buf._AtIndex(0)) {
		return parseNumberVal(obj, buf)
	}

	return fmt.Errorf("unexpected character with code %d at position %d", buf._AtIndex(0), buf._GetOffset())
}

func parseStringVal(obj *_TscObj, buf *_Buffer) error {
	if !buf._CanAccessAtIndex(0) {
		return fmt.Errorf("expected \" on the position %d", buf)
	}

	if !buf._CmpCharAtIndex(0, '"') {
		return fmt.Errorf("unexpected %s at position %d", string(buf._AtIndex(0)), buf._GetOffset())
	}

	buf._AdvanceAt(1)

	var (
		i int

		allocStr     []byte
		skippedBytes int
	)

	skippedBytes = 0

	for i = 0; buf._CanAccessAtIndex(uint(i)) && buf._AtIndex(uint(i)) != '"' && buf._AtIndex(uint(i)) > 32; i++ {
		if buf._AtIndex(uint(i)) == '\\' {
			skippedBytes++
			i++
		}
	}

	if !(buf._CanAccessAtIndex(uint(i)) && buf._AtIndex(uint(i)) == '"') {
		return fmt.Errorf("string ended unexpectedly, unexpected character %d, start of string at position %d", buf._AtIndex(uint(i)), buf._GetOffset())
	}

	allocStr = make([]byte, i-skippedBytes)

	var (
		aIndex int
		mIndex int
	)

	mIndex = 0

	// copy
	for aIndex = 0; mIndex < i; aIndex++ {
		if buf._AtIndex(uint(mIndex)) != '\\' {
			allocStr[aIndex] = buf._AtIndex(uint(mIndex))
			mIndex++
		} else {
			mIndex++

			switch buf._AtIndex(uint(mIndex)) {
			case 'n':
				allocStr[aIndex] = '\n'
				break
			case 't':
				allocStr[aIndex] = '\t'
				break
			case 'r':
				allocStr[aIndex] = '\r'
				break
			case 'b':
				allocStr[aIndex] = '\b'
				break
			case '"':
			case '\\':
				allocStr[aIndex] = buf._AtIndex(uint(mIndex))
				break
			default:
				return fmt.Errorf("unreachable escaped char %s", string(buf._AtIndex(uint(mIndex))))
			}

			mIndex++
		}
	}

	obj.T = _ObjStringT
	obj.StrVal = allocStr

	buf._AdvanceAt(uint(i + 1))

	return nil
}

func parseNumberVal(obj *_TscObj, buf *_Buffer) error {
	if !buf._CanAccessAtIndex(0) {
		return fmt.Errorf("unexpected end")
	}

	var (
		i int

		intVal int64

		err error
	)

	for i = 0; buf._CanAccessAtIndex(uint(i)); i++ {
		switch buf._AtIndex(uint(i)) {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
		case '-':
			break
		default:
			goto loop_end
		}
	}
loop_end:
	if intVal, err = strconv.ParseInt(string(buf._MakeSlice(0, uint(i))), 10, 32); err != nil {
		return fmt.Errorf("parse number err at position %d, err: %s", buf._GetOffset(), err.Error())
	}

	obj.T = _ObjIntT
	obj.Int = int(intVal)

	buf._AdvanceAt(uint(i))

	return nil
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

// available only number, uppercase character, lowercase character
func isKeyCharacter(char byte) bool {
	if (char >= 97 && char <= 122) || (char >= 65 && char <= 90) || (char >= 48 && char <= 57) {
		return true
	}

	return false
}

func isANumber(char byte) bool {
	return char >= '0' && char <= '9'
}
