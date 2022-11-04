package pgdb

import (
	"context"
	"database/sql"
)

// hookableConn wraps a Conn alongside a list of Hooks to be called.
// Each hook's Before method is called in the order they are specified
// before passing the call onto the wrapped Conn.
// The functions returned by the hooks' Before methods are called
// after the response is received from the db in FILO order.
type hookableConn struct {
	db    Conn
	hooks []Hook
}

// Runs any registered hooks followed by Conn.ExecContext
func (c *hookableConn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	defer c.runHooks(query, args...)()
	return c.db.ExecContext(ctx, query, args...)
}

// Runs any registered hooks followed by Conn.QueryContext
func (c *hookableConn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	defer c.runHooks(query, args...)()
	return c.db.QueryContext(ctx, query, args...)
}

// Runs any registered hooks followed by Conn.QueryRowContext
func (c *hookableConn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	defer c.runHooks(query, args...)()
	return c.db.QueryRowContext(ctx, query, args...)
}

// Runs any registered hooks followed by Conn.Exec
func (c *hookableConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	defer c.runHooks(query, args...)()
	return c.db.Exec(query, args...)
}

// Runs any registered hooks followed by Conn.Query
func (c *hookableConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	defer c.runHooks(query, args...)()
	return c.db.Query(query, args...)
}

// Runs any registered hooks followed by Conn.QueryRow
func (c *hookableConn) QueryRow(query string, args ...interface{}) *sql.Row {
	defer c.runHooks(query, args...)()
	return c.db.QueryRow(query, args...)
}

func (c *hookableConn) runHooks(query string, args ...interface{}) func() {
	var afterFns []func()
	for _, hook := range c.hooks {
		fn := hook.Before(query, args...)
		if fn != nil {
			afterFns = append(afterFns, fn)
		}
	}
	return func() {
		for i := len(afterFns) - 1; i >= 0; i-- {
			afterFns[i]()
		}
	}
}
