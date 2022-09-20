package pager

import (
	"bytes"
	"context"
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
)

const (
	MaxRows int = 20
)

var pool *sync.Pool = &sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

type MetaInfo struct {
	Total  int
	Limit  int
	Offset int
}

func NewPageMeta() MetaInfo {
	return MetaInfo{
		Total:  0,
		Limit:  0,
		Offset: 0,
	}
}

func (m *MetaInfo) NumberOfPages() int {
	if m.Limit >= m.Total {
		return 1
	}

	if (m.Total % m.Limit) > 0 {
		return (m.Total / m.Limit) + 1
	}

	return m.Total / m.Limit

}

func (m *MetaInfo) CurrentPage() int {
	if m.Offset >= m.Total {
		return 0
	}

	return (m.Offset / m.Limit) + 1
}

func FetchPage(ctx context.Context, conn *sqlx.DB, m *MetaInfo, dest interface{}, table string, conds []string) error {
	var (
		filterBuf *bytes.Buffer
		queryBuf  *bytes.Buffer

		filterStr string

		err error
	)

	// query
	queryBuf = pool.Get().(*bytes.Buffer)
	defer queryBuf.Reset()
	defer pool.Put(queryBuf)

	// filter
	filterBuf = pool.Get().(*bytes.Buffer)

	_ParseAndBuildConds(filterBuf, conds, &m.Limit, &m.Offset)

	if m.Limit == 0 {
		m.Limit = MaxRows
	}

	filterBuf.WriteByte(' ')
	filterBuf.WriteString("limit ")
	filterBuf.WriteString(strconv.Itoa(m.Limit))

	filterBuf.WriteByte(' ')
	filterBuf.WriteString("offset ")
	filterBuf.WriteString(strconv.Itoa(m.Offset))

	filterStr = filterBuf.String()
	filterBuf.Reset()
	pool.Put(filterBuf)

	// prepare query to fetch count of rows
	buildQuery(queryBuf, "select count(*) as totals from", table, filterStr)

	// exec query
	if err = conn.QueryRowContext(ctx, queryBuf.String()).Scan(&m.Total); err != nil {
		return err
	}

	// reset previous query row
	queryBuf.Reset()

	// prepare query to fetch resources
	buildQuery(queryBuf, "select * from", table, filterStr)

	if err = conn.SelectContext(ctx, dest, queryBuf.String()); err != nil {
		return err
	}

	return nil
}

func buildQuery(buf *bytes.Buffer, sPart string, table string, filterPart string) {
	buf.WriteString(sPart)
	buf.WriteByte(' ')
	buf.WriteString(table)
	buf.WriteByte(' ')
	buf.WriteString(filterPart)
}
