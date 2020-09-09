// Copyright (c) 2017-2020 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
)

type paramKey string

const (
	// MultiStatementCount controls the number of queries to execute in a single API call
	MultiStatementCount paramKey = "MULTI_STATEMENT_COUNT"
	// AsyncMode controls
	AsyncMode paramKey = "ASYNC_MODE_QUERY"
	// QueryIDChan controls
	QueryIDChan paramKey = "QUERY_ID_CHAN"
	// ResumeQueryID controls
	ResumeQueryID paramKey = "RESUME_QUERY_ID"
	// IsInternal controls
	IsInternal paramKey = "INTERNAL_QUERY"
)

type snowflakeStmt struct {
	sc    *snowflakeConn
	query string
}

func (stmt *snowflakeStmt) Close() error {
	glog.V(2).Infoln("Stmt.Close")
	// noop
	return nil
}

func (stmt *snowflakeStmt) NumInput() int {
	glog.V(2).Infoln("Stmt.NumInput")
	// Go Snowflake doesn't know the number of binding parameters.
	return -1
}

func (stmt *snowflakeStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	glog.V(2).Infoln("Stmt.ExecContext")
	return stmt.sc.ExecContext(ctx, stmt.query, args)
}

func (stmt *snowflakeStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	glog.V(2).Infoln("Stmt.QueryContext")
	return stmt.sc.QueryContext(ctx, stmt.query, args)
}

func (stmt *snowflakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	glog.V(2).Infoln("Stmt.Exec")
	return stmt.sc.Exec(stmt.query, args)
}

func (stmt *snowflakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	glog.V(2).Infoln("Stmt.Query")
	return stmt.sc.Query(stmt.query, args)
}

// WithMultiStatement returns a context that allows the user to execute the desired number of sql queries in one query
func WithMultiStatement(ctx context.Context, num int) (context.Context, error) {
	return context.WithValue(ctx, MultiStatementCount, num), nil
}

// WithAsyncMode returns a context
func WithAsyncMode(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, AsyncMode, true), nil
}

// WithInternal returns a context
func WithInternal(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, IsInternal, true), nil
}
