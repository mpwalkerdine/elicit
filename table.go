package elicit

// Table holds test data from the spec
type Table struct {
	Columns []string
	Rows    []map[string]string
}

func makeTable(rows [][]string) Table {
	t := Table{
		Columns: rows[0],
	}

	cc := len(t.Columns)

	for _, row := range rows[1:] {
		m := make(map[string]string, cc)
		for c := 0; c < cc; c++ {
			m[t.Columns[c]] = row[c]
		}
		t.Rows = append(t.Rows, m)
	}

	return t
}
