// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package databasesql

import (
	"context"
	"database/sql"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"log"
	"strings"
)

var databaseSqlInstrumenter = BuildDatabaseSqlOtelInstrumenter()

var dbSqlEnabler = instrumenter.NewDefaultInstrumentEnabler()

// ------ Open ------

// func Open(driverName, dataSourceName string) (*DB, error)
func beforeOpenInstrumentation(call api.CallContext, driverName, dataSourceName string) {
	if !dbSqlEnabler.Enable() { // Always return true?
		return
	}
	addr, err := parseDSN(driverName, dataSourceName)
	if err != nil {
		log.Printf("failed to parse dsn: %v", err)
	}
	call.SetData(map[string]string{
		"endpoint": addr,
		"driver":   driverName,
		"dsn":      dataSourceName,
	})
}

// func Open(driverName, dataSourceName string) (*DB, error)
func afterOpenInstrumentation(call api.CallContext, db *sql.DB, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	endpoint, ok := data["endpoint"]
	if ok {
		db.Endpoint = endpoint
	}
	driver, ok := data["driver"]
	if ok {
		db.DriverName = driver
	}
	dsn, ok := data["dsn"]
	if ok {
		db.DSN = dsn
	}
}

// ------ Ping ------

// func (db *DB) PingContext(ctx context.Context) error
func beforePingContextInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "ping", "ping", db.Endpoint, db.DriverName, db.DSN)
}

func afterPingContextInstrumentation(call api.CallContext, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

// func (db *DB) PingContext(ctx context.Context) error
func beforePrepareContextInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context, query string) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": db.Endpoint,
		"sql":      query,
		"driver":   db.DriverName,
		"dsn":      db.DSN,
	})
}

func afterPrepareContextInstrumentation(call api.CallContext, stmt *sql.Stmt, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeExecContextInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, db.Endpoint, db.DriverName, db.DSN, args...)
}

func afterExecContextInstrumentation(call api.CallContext, result sql.Result, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

// ------ Query ------
// func (db *DB) QueryContext(ctx context.Context, query string, args ...any) 
func beforeQueryContextInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, db.Endpoint, db.DriverName, db.DSN, args...)
}

// func (db *DB) QueryContext(ctx context.Context, query string, args ...any) 
func afterQueryContextInstrumentation(call api.CallContext, rows *sql.Rows, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

// ------ TX ------

func beforeTxInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context, opts *sql.TxOptions) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "begin", "START TRANSACTION", db.Endpoint, db.DriverName, db.DSN)
}

func afterTxInstrumentation(call api.CallContext, tx *sql.Tx, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	callData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	dbRequest, ok := callData["dbRequest"].(databaseSqlRequest)
	if !ok {
		return
	}
	tx.Endpoint = dbRequest.endpoint
	tx.DriverName = dbRequest.driverName
	tx.DSN = dbRequest.dsn
	instrumentEnd(call, err)
}

func beforeConnInstrumentation(call api.CallContext, db *sql.DB, ctx context.Context) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": db.Endpoint,
		"driver":   db.DriverName,
		"dsn":      db.DSN,
	})
}

func afterConnInstrumentation(call api.CallContext, conn *sql.Conn, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	endpoint, ok := data["endpoint"]
	if ok {
		conn.Endpoint = endpoint
	}
	driverName, ok := data["driver"]
	if ok {
		conn.DriverName = driverName
	}
	dsn, ok := data["dsn"]
	if ok {
		conn.DSN = dsn
	}
}

func beforeConnPingContextInstrumentation(call api.CallContext, conn *sql.Conn, ctx context.Context) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "ping", "ping", conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnPingContextInstrumentation(call api.CallContext, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeConnPrepareContextInstrumentation(call api.CallContext, conn *sql.Conn, ctx context.Context, query string) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": conn.Endpoint,
		"sql":      query,
		"driver":   conn.DriverName,
		"dsn":      conn.DSN,
	})
}

func afterConnPrepareContextInstrumentation(call api.CallContext, stmt *sql.Stmt, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeConnExecContextInstrumentation(call api.CallContext, conn *sql.Conn, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, conn.Endpoint, conn.DriverName, conn.DSN, args...)
}

func afterConnExecContextInstrumentation(call api.CallContext, result sql.Result, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeConnQueryContextInstrumentation(call api.CallContext, conn *sql.Conn, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, conn.Endpoint, conn.DriverName, conn.DSN, args...)
}

func afterConnQueryContextInstrumentation(call api.CallContext, rows *sql.Rows, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeConnTxInstrumentation(call api.CallContext, conn *sql.Conn, ctx context.Context, opts *sql.TxOptions) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "start", "START TRANSACTION", conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnTxInstrumentation(call api.CallContext, tx *sql.Tx, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeTxPrepareContextInstrumentation(call api.CallContext, tx *sql.Tx, ctx context.Context, query string) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": tx.Endpoint,
		"sql":      query,
		"driver":   tx.DriverName,
		"dsn":      tx.DSN,
	})
}

func afterTxPrepareContextInstrumentation(call api.CallContext, stmt *sql.Stmt, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeTxStmtContextInstrumentation(call api.CallContext, tx *sql.Tx, ctx context.Context, stmt *sql.Stmt) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": stmt.Data["endpoint"],
		"driver":   stmt.Data["driver"],
		"dsn":      stmt.DSN,
	})
}

func afterTxStmtContextInstrumentation(call api.CallContext, stmt *sql.Stmt) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{}
	endpoint, ok := data["endpoint"]
	if ok {
		stmt.Data["endpoint"] = endpoint
	}
	driverName, ok := data["driver"]
	if ok {
		stmt.Data["driver"] = driverName
	}
	dsn, ok := data["dsn"]
	if ok {
		stmt.Data["dsn"] = dsn
	}
}

func beforeTxExecContextInstrumentation(call api.CallContext, tx *sql.Tx, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, tx.Endpoint, tx.DriverName, tx.DSN, args...)
}

func afterTxExecContextInstrumentation(call api.CallContext, result sql.Result, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeTxQueryContextInstrumentation(call api.CallContext, tx *sql.Tx, ctx context.Context, query string, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, tx.Endpoint, tx.DriverName, tx.DSN, args...)
}

func afterTxQueryContextInstrumentation(call api.CallContext, rows *sql.Rows, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeTxCommitInstrumentation(call api.CallContext, tx *sql.Tx) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	instrumentStart(call, context.Background(), "commit", "COMMIT", tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxCommitInstrumentation(call api.CallContext, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeTxRollbackInstrumentation(call api.CallContext, tx *sql.Tx) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if tx == nil {
		return
	}
	instrumentStart(call, context.Background(), "rollback", "ROLLBACK", tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxRollbackInstrumentation(call api.CallContext, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeStmtExecContextInstrumentation(call api.CallContext, stmt *sql.Stmt, ctx context.Context, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	sql, endpoint, driverName, dsn := "", "", "", ""
	if stmt.Data != nil {
		sql, endpoint, driverName, dsn = stmt.Data["sql"], stmt.Data["endpoint"], stmt.Data["driver"], stmt.DSN
	}
	instrumentStart(call, ctx, "exec", sql, endpoint, driverName, dsn, args...)
}

func afterStmtExecContextInstrumentation(call api.CallContext, result sql.Result, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func beforeStmtQueryContextInstrumentation(call api.CallContext, stmt *sql.Stmt, ctx context.Context, args ...any) {
	if !dbSqlEnabler.Enable() {
		return
	}
	if stmt == nil {
		return
	}
	sql, endpoint, driverName, dsn := "", "", "", ""
	if stmt.Data != nil {
		sql, endpoint, driverName, dsn = stmt.Data["sql"], stmt.Data["endpoint"], stmt.Data["driver"], stmt.DSN
	}
	instrumentStart(call, ctx, "query", sql, endpoint, driverName, dsn, args...)
}

func afterStmtQueryContextInstrumentation(call api.CallContext, rows *sql.Rows, err error) {
	if !dbSqlEnabler.Enable() {
		return
	}
	instrumentEnd(call, err)
}

func instrumentStart(call api.CallContext, ctx context.Context, spanName, query, endpoint, driverName, dsn string, args ...any) {
	req := databaseSqlRequest{
		opType:     calOp(query),
		sql:        query,
		endpoint:   endpoint,
		driverName: driverName,
		dsn:        dsn,
		params:     args,
	}
	newCtx := databaseSqlInstrumenter.Start(ctx, req)
	call.SetData(map[string]interface{}{
		"dbRequest": req,
		"newCtx":    newCtx,
	})
}

func instrumentEnd(call api.CallContext, err error) {
	callData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	dbRequest, ok := callData["dbRequest"].(databaseSqlRequest)
	if !ok {
		return
	}
	newCtx, ok := callData["newCtx"].(context.Context)
	if !ok {
		return
	}
	databaseSqlInstrumenter.End(newCtx, dbRequest, nil, err)
}

func calOp(sql string) string {
	sqls := strings.Split(sql, " ")
	var op string
	if len(sqls) > 0 {
		op = sqls[0]
	}
	return op
}
