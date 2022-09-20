package pager

import (
	"bytes"
	"strconv"
)

var (
	sortParam    = []byte("sort")
	limitParam   = []byte("limit")
	offsetParam  = []byte("offset")
	orderByParam = []byte("order")

	desc = []byte("desc")
	asc  = []byte("asc")
)

/*
IN EXAMPLE: ...?sort=updatedAt&order=desc&search=id&offset=0&status=1&text#1


IN

ConditionsList
[
	"sort=updatedAt",
	...,
]

"sort=updatedAt"

'sort' - block
'updatedAt' - value
'=' -  separtes black and


[10, 2, 3, 20, 1, 5, 6]

OUT

ORDER BY created_at desc

*/

func _ParseAndBuildConds(buf *bytes.Buffer, conds []string, limit *int, offset *int) {
	var (
		c string
		i int
		n int

		bv    []byte
		param []byte

		sep byte
	)

	sep = 0

	//sort(conds)

	for _, c = range conds {
		bv = []byte(c)
		n = len(bv)

		for i = 0; i < n && bv[i] != '='; i++ {
			// spin
		}

		if i == n || (i+1) == n {
			continue
		}

		param = bv[0:i]

		// sort
		if bytes.Compare(sortParam, param) == 0 {
			buf.WriteByte(sep)

			if parseSort(buf, bv[i+1:]) {
				sep = ' '
			}
		}

		// limit
		if bytes.Compare(limitParam, param) == 0 {
			*limit = forceParseInt(bv[i+1:])
		}

		// offset
		if bytes.Compare(offsetParam, param) == 0 {
			*offset = forceParseInt(bv[i+1:])
		}
	}

}

func parseSort(buf *bytes.Buffer, val []byte) bool {
	var (
		i            int
		j            int
		h            int
		k            int
		n            int
		stmIsWritten bool
		resetF       bool
		sep          byte
		order        []byte
	)

	// for writing the order word only when we write some a handled condition to buffer
	stmIsWritten = false

	sep = 0

	j = 0
	k = 0
	h = 0

	resetF = false

	order = asc

	n = len(val)
loop:
	for i < n {
		if resetF {
			sep = ','
			order = asc

			j = i
			k = i

			resetF = false

		}

		switch val[i] {
		case '!':
			for h = i + 1; h < n && val[h] != ','; h++ {
				// spin
			}

			if bytes.Compare(desc, val[i+1:h]) == 0 {
				order = desc
			}

			if bytes.Compare(asc, val[i+1:h]) == 0 {
				order = asc
			}

			i = h

			break
			// created_at,
		case ',':
			i++

			if j < k {
				goto write
			}

			break
		default:
			k++
			i++
		}
	}

write:
	if j < k {
		if !stmIsWritten {
			buf.WriteString("order by ")
			stmIsWritten = true
		}

		buf.WriteByte(sep)
		buf.Write(val[j:k])
		buf.WriteByte(' ')
		buf.Write(order)

		if i < n {
			resetF = true
			goto loop
		}
	}

	return stmIsWritten == true
}

func forceParseInt(a []byte) int {
	var (
		r   int64
		err error
	)

	if r, err = strconv.ParseInt(string(a), 10, 32); err != nil {
		return 0
	}

	return int(r)
}

func parseOffset() int {
	return 0
}

func parseLimit() int {

	return 0
}

func canRead(a []byte, n int) bool {
	return n >= len(a)
}

func sort(conds []string) {
	// not implemented
	// it is important for right order so that there is not syntax sql error, but we ignore this moment for now
}
