package handler

import (
	"testing"

	"github.com/lighttiger2505/sqls/internal/config"
	"github.com/lighttiger2505/sqls/internal/database"
	"github.com/lighttiger2505/sqls/internal/lsp"
)

type completionTestCase struct {
	name  string
	input string
	line  int
	col   int
	want  []string
}

var statementCase = []completionTestCase{
	{
		name:  "columns on multiple statement forcused first",
		input: "SELECT c. FROM city as c;SELECT c. FROM country as c;",
		line:  0,
		col:   9,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "columns with multiple statement forcused second",
		input: "SELECT c. FROM city as c;SELECT c. FROM country as c;",
		line:  0,
		col:   34,
		want: []string{
			"Code",
			"Name",
			"CountryCode",
			"Region",
			"SurfaceArea",
			"IndepYear",
			"LifeExpectancy",
			"GNP",
			"GNPOld",
			"LocalName",
			"GovernmentForm",
			"HeadOfState",
			"Capital",
			"Code2",
		},
	},
}

var selectExprCase = []completionTestCase{
	{
		name:  "table columns",
		input: "select  from city",
		line:  0,
		col:   7,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "filterd table columns",
		input: "select Cou from city",
		line:  0,
		col:   10,
		want: []string{
			"CountryCode",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "columns of table specifies database",
		input: "select  from world.city",
		line:  0,
		col:   7,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "columns of aliased table",
		input: "select  from city as c",
		line:  0,
		col:   7,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"c",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "columns of aliased table specifies database",
		input: "select  from world.city as c",
		line:  0,
		col:   7,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"c",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "columns of aliased table without as",
		input: "select  from city c",
		line:  0,
		col:   7,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"c",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "columns of before period table",
		input: "select c. from city as c",
		line:  0,
		col:   9,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "filterd columns of before period table",
		input: "select c.C from city as c",
		line:  0,
		col:   10,
		want: []string{
			"CountryCode",
		},
	},
	{
		name:  "identifier list",
		input: "select id,  from city",
		line:  0,
		col:   11,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "filterd identifier list",
		input: "select id, cou from city",
		line:  0,
		col:   14,
		want: []string{
			"CountryCode",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "comparison",
		input: "select 1 = cou from city",
		line:  0,
		col:   14,
		want: []string{
			"CountryCode",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "operator",
		input: "select 1 + cou from city",
		line:  0,
		col:   14,
		want: []string{
			"CountryCode",
			"country",
			"countrylanguage",
		},
	},
}

var tableReferenceCase = []completionTestCase{
	{
		name:  "from tables",
		input: "select CountryCode from ",
		line:  0,
		col:   24,
		want: []string{
			"city",
			"country",
			"countrylanguage",
			"information_schema",
			"mysql",
			"performance_schema",
			"sys",
			"world",
		},
	},
	{
		name:  "from filterd tables",
		input: "select CountryCode from co",
		line:  0,
		col:   26,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "join tables",
		input: "select CountryCode from city join ",
		line:  0,
		col:   34,
		want: []string{
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "join filterd tables",
		input: "select CountryCode from city join co",
		line:  0,
		col:   36,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "insert tables",
		input: "INSERT INTO ",
		line:  0,
		col:   12,
		want: []string{
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "insert filterd tables",
		input: "INSERT INTO co",
		line:  0,
		col:   12,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "insert columns",
		input: "INSERT INTO city (",
		line:  0,
		col:   18,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "insert filterd columns",
		input: "INSERT INTO city (cou",
		line:  0,
		col:   21,
		want: []string{
			"CountryCode",
		},
	},
	{
		name:  "insert identifier list",
		input: "INSERT INTO city (id, cou",
		line:  0,
		col:   25,
		want: []string{
			"CountryCode",
		},
	},
	{
		name:  "update tables",
		input: "UPDATE ",
		line:  0,
		col:   7,
		want: []string{
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "update filterd tables",
		input: "UPDATE co",
		line:  0,
		col:   9,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "update columns",
		input: "UPDATE city SET ",
		line:  0,
		col:   16,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "update filterd columns",
		input: "UPDATE city SET cou",
		line:  0,
		col:   19,
		want: []string{
			"CountryCode",
		},
	},
	{
		name:  "update identiger list",
		input: "UPDATE city SET CountryCode=12, Na",
		line:  0,
		col:   34,
		want: []string{
			"Name",
		},
	},
	{
		name:  "delete tables",
		input: "DELETE FROM ",
		line:  0,
		col:   12,
		want: []string{
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "delete filterd tables",
		input: "DELETE FROM co",
		line:  0,
		col:   14,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
}

var whereCondition = []completionTestCase{
	{
		name:  "where columns",
		input: "select * from city where ",
		line:  0,
		col:   25,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "where columns of specified table",
		input: "select * from city where city.",
		line:  0,
		col:   30,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "where columns in left of comparison",
		input: "select * from city where  = ID",
		line:  0,
		col:   25,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "where columns in right of comparison",
		input: "select * from city where ID = ",
		line:  0,
		col:   30,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "where columns of specified table in left of comparison",
		input: "select * from city where city. = city.ID",
		line:  0,
		col:   30,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "where columns of specified table in right of comparison",
		input: "select * from city where city.ID = city.",
		line:  0,
		col:   40,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
		},
	},
	{
		name:  "join on columns",
		input: "select * from city left join country on ",
		line:  0,
		col:   40,
		want: []string{
			"Code",
			"Name",
			"CountryCode",
			"Continent",
			"Region",
			"SurfaceArea",
			"IndepYear",
			"LifeExpectancy",
			"GNP",
			"GNPOld",
			"LocalName",
			"GovernmentForm",
			"HeadOfState",
			"Capital",
			"Code2",
		},
	},
	{
		name:  "join on filterd columns",
		input: "select * from city left join country on co",
		line:  0,
		col:   52,
		want: []string{
			"Code",
			"Continent",
			"Code2",
		},
	},
}

var colNameCase = []completionTestCase{
	{
		name:  "ORDER BY columns",
		input: "SELECT ID, Name FROM city ORDER BY ",
		line:  0,
		col:   35,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "GROUP BY columns",
		input: "SELECT CountryCode, COUNT(*) FROM city GROUP BY ",
		line:  0,
		col:   48,
		want: []string{
			"ID",
			"Name",
			"CountryCode",
			"District",
			"Population",
			"city",
			"country",
			"countrylanguage",
		},
	},
}

var caseValueCase = []completionTestCase{
	{
		name:  "select case after case",
		input: "SELECT `Language`, CASE Is WHEN 'T' THEN 'official' WHEN 'F' THEN 'unofficial' END AS is_official FROM countrylanguage;",
		line:  0,
		col:   26,
		want: []string{
			"IsOfficial",
		},
	},
	{
		name:  "select case after when",
		input: "SELECT `Language`, CASE IsOfficial WHEN Is THEN 'official' WHEN 'F' THEN 'unofficial' END AS is_official FROM countrylanguage;",
		line:  0,
		col:   42,
		want: []string{
			"IsOfficial",
		},
	},
	{
		name:  "select case after then",
		input: "SELECT `Language`, CASE IsOfficial WHEN 'T' THEN Is WHEN 'F' THEN 'unofficial' END AS is_official FROM countrylanguage;",
		line:  0,
		col:   51,
		want: []string{
			"IsOfficial",
		},
	},
	{
		name:  "select case identifier list",
		input: "SELECT `Language`, CASE IsOfficial WHEN 'T' THEN Is WHEN 'F' THEN 'unofficial' END AS is_official, P FROM countrylanguage;",
		line:  0,
		col:   100,
		want: []string{
			"Percentage",
		},
	},
}

var subQueryCase = []completionTestCase{
	{
		name:  "in subquery table columns",
		input: "SELECT * FROM (SELECT Cou FROM city)",
		line:  0,
		col:   25,
		want: []string{
			"CountryCode",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "in subquery table references",
		input: "SELECT * FROM (SELECT * FROM ",
		line:  0,
		col:   29,
		want: []string{
			"city",
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "in subquery filterd table references",
		input: "SELECT * FROM (SELECT * FROM co",
		line:  0,
		col:   29,
		want: []string{
			"country",
			"countrylanguage",
		},
	},
	{
		name:  "subquery columns",
		input: "SELECT  FROM (SELECT ID as city_id, Name as city_name FROM city) as t",
		line:  0,
		col:   7,
		want: []string{
			"t",
			"city_id",
			"city_name",
		},
	},
	{
		name:  "subquery columns",
		input: "SELECT t. FROM (SELECT ID as city_id, Name as city_name FROM city) as t",
		line:  0,
		col:   9,
		want: []string{
			"city_id",
			"city_name",
		},
	},
	{
		name:  "columns of multiple subquery",
		input: "SELECT  FROM (SELECT Name as city_name FROM city) AS sub1, (SELECT LocalName as country_name FROM country) AS sub2 limit 1",
		line:  0,
		col:   7,
		want: []string{
			"sub1",
			"sub2",
			"city_name",
			"country_name",
		},
	},
}

func TestComplete(t *testing.T) {
	tx := newTestContext()
	tx.initServer(t)
	defer tx.tearDown()

	didChangeConfigurationParams := lsp.DidChangeConfigurationParams{
		Settings: struct {
			SQLS *config.Config "json:\"sqls\""
		}{
			SQLS: &config.Config{
				Connections: []*database.DBConfig{
					{
						Driver:         "mock",
						DataSourceName: "",
					},
				},
			},
		},
	}
	if err := tx.conn.Call(tx.ctx, "workspace/didChangeConfiguration", didChangeConfigurationParams, nil); err != nil {
		t.Fatal("conn.Call workspace/didChangeConfiguration:", err)
	}

	uri := "file:///Users/octref/Code/css-test/test.sql"
	testcaseMap := map[string][]completionTestCase{
		"statement":       statementCase,
		"select expr":     selectExprCase,
		"table reference": tableReferenceCase,
		"col name":        colNameCase,
		"case value":      caseValueCase,
		"subquery":        subQueryCase,
	}

	for k, v := range testcaseMap {
		for _, tt := range v {
			t.Run(k+" "+tt.name, func(t *testing.T) {
				// Open dummy file
				didOpenParams := lsp.DidOpenTextDocumentParams{
					TextDocument: lsp.TextDocumentItem{
						URI:        uri,
						LanguageID: "sql",
						Version:    0,
						Text:       tt.input,
					},
				}
				if err := tx.conn.Call(tx.ctx, "textDocument/didOpen", didOpenParams, nil); err != nil {
					t.Fatal("conn.Call textDocument/didOpen:", err)
				}
				tx.testFile(t, didOpenParams.TextDocument.URI, didOpenParams.TextDocument.Text)
				// Create completion params
				commpletionParams := lsp.CompletionParams{
					TextDocumentPositionParams: lsp.TextDocumentPositionParams{
						TextDocument: lsp.TextDocumentIdentifier{
							URI: uri,
						},
						Position: lsp.Position{
							Line:      tt.line,
							Character: tt.col,
						},
					},
					CompletionContext: lsp.CompletionContext{
						TriggerKind:      0,
						TriggerCharacter: nil,
					},
				}

				var got []lsp.CompletionItem
				if err := tx.conn.Call(tx.ctx, "textDocument/completion", commpletionParams, &got); err != nil {
					t.Fatal("conn.Call textDocument/completion:", err)
				}
				testCompletionItem(t, tt.want, got)
			})
		}
	}

}

func testCompletionItem(t *testing.T, expectLabels []string, gotItems []lsp.CompletionItem) {
	t.Helper()

	itemMap := map[string]struct{}{}
	for _, item := range gotItems {
		itemMap[item.Label] = struct{}{}
	}

	for _, el := range expectLabels {
		_, ok := itemMap[el]
		if !ok {
			t.Errorf("not included label, expect %q", el)
		}
	}
}
