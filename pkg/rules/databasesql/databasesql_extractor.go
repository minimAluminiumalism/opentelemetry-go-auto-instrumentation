// Copyright (c) 2025 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package databasesql

import (
	"fmt"
	"log"

	"github.com/xwb1989/sqlparser"
)

func extractSQLMetadata(request databaseSqlRequest) {
	sql := request.sql

	if sqlCache.Contains(sql) {
		return
	}

	collection := extractCollection(sql)
	sqlMeta := SQLMeta{
		stmt:       request.sql,
		operation:  request.opType,
		collection: collection,
	}

	sqlCache.Add(sql, sqlMeta)
}

func getCollection(sql string) string {
	if meta, found := sqlCache.Get(sql); found {
		return meta.collection
	}
	// Attempt to retrieve the collection again.
	return extractCollection(sql)
}

func getParams(sql string) []any {
	meta, found := sqlCache.Get(sql)
	if found && len(meta.params) > 0 {
		return meta.params
	}

	// db params did not been retrieved in `extractSQLMetadata()`, so parse it and update the db meta here.
	paramsMap, _ := extractSQLParams(sql)
	params := []any{}
	for _, v := range paramsMap {
		params = append(params, v)
	}
	if found {
		// `meta` was found already, only need to update it's params.
		updatedMeta := meta
		updatedMeta.params = params
		sqlCache.Add(sql, updatedMeta)
	}
	return params
}

func extractCollection(query string) string {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}
	// Only support DML currently
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		return getTableName(stmt.From)
	case *sqlparser.Update:
		return getTableName(stmt.TableExprs)
	case *sqlparser.Insert:
		return stmt.Table.Name.String()
	case *sqlparser.Delete:
		return getTableName(stmt.TableExprs)
	default:
		return ""
	}
}

func getTableName(node sqlparser.SQLNode) string {
	switch n := node.(type) {
	case sqlparser.TableName:
		return n.Name.String()
	case sqlparser.TableExprs:
		for _, expr := range n {
			aliasedExpr, ok := expr.(*sqlparser.AliasedTableExpr)
			if !ok {
				continue
			}
			tableName, ok := aliasedExpr.Expr.(sqlparser.TableName)
			if ok {
				return tableName.Name.String()
			}
		}
	}
	return ""
}

// Extract SQL parameters
func extractSQLParams(query string) (map[string]string, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		log.Printf("failed to fetch sql params: %v", err)
		return nil, err
	}

	values := make(map[string]string)
	// Only support DML currently
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		if stmt.Where != nil {
			extractConditions(stmt.Where.Expr, values)
		}
	case *sqlparser.Update:
		for _, expr := range stmt.Exprs {
			if sqlVal, ok := expr.Expr.(*sqlparser.SQLVal); ok && sqlVal.Type == sqlparser.StrVal {
				values[expr.Name.Name.String()] = string(sqlVal.Val)
			}
		}
		if stmt.Where != nil {
			extractConditions(stmt.Where.Expr, values)
		}
	case *sqlparser.Delete:
		if stmt.Where != nil {
			extractConditions(stmt.Where.Expr, values)
		}
	case *sqlparser.Insert:
		columns := make([]string, 0, len(stmt.Columns))
		for _, col := range stmt.Columns {
			columns = append(columns, col.String())
		}

		rows, ok := stmt.Rows.(sqlparser.Values)
		if ok {
			for i, row := range rows {
				rowSuffix := ""
				if len(rows) > 1 {
					rowSuffix = fmt.Sprintf("_row%d", i+1)
				}

				for j, val := range row {
					if j >= len(columns) {
						continue
					}

					if sqlVal, ok := val.(*sqlparser.SQLVal); ok && sqlVal.Type == sqlparser.StrVal {
						colName := columns[j]
						if rowSuffix != "" {
							colName += rowSuffix
						}
						values[colName] = string(sqlVal.Val)
					}
				}
			}
		}
	}

	return values, nil
}

func extractConditions(expr sqlparser.Expr, values map[string]string) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if expr.Operator == "=" {
			colName, ok := expr.Left.(*sqlparser.ColName)
			if !ok {
				return
			}

			valExpr, ok := expr.Right.(*sqlparser.SQLVal)
			if !ok {
				return
			}

			if valExpr.Type == sqlparser.StrVal {
				values[colName.Name.String()] = string(valExpr.Val)
			}
		}
	case *sqlparser.AndExpr:
		extractConditions(expr.Left, values)
		extractConditions(expr.Right, values)
	case *sqlparser.OrExpr:
		extractConditions(expr.Left, values)
		extractConditions(expr.Right, values)
	}
}
