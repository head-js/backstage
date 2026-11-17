//go:build ignore

package main

// Yaegi loads this file after the types and HTTP client, before the translator.
// Direct database access to the manually added tasks.ext_* columns (see README).

import (
	storefmt "fmt"
	storestrconv "strconv"

	"code.vikunja.io/api/pkg/db"
)

// taskExt carries the plugin-owned extension columns of one task row. Adding a
// new ext column means: ALTER TABLE, a field here, and the SELECT lists below.
// NULL columns read back as empty strings.
type taskExt struct {
	Iid     string // tasks.ext_iid
	Session string // tasks.ext_session
}

// findTaskExtByIid resolves a stored VisitNode Id (tasks.ext_iid) to the
// native task id and the row's extension columns, in one query.
func findTaskExtByIid(visitNodeId string) (taskId int64, ext taskExt, found bool, err error) {
	s := db.NewAutocommitSession()
	defer s.Close()
	rows, err := s.Query(`SELECT id, ext_iid, ext_session FROM tasks WHERE ext_iid = ? AND deleted_at IS NULL`, visitNodeId)
	if err != nil {
		return 0, taskExt{}, false, err
	}
	if len(rows) == 0 {
		return 0, taskExt{}, false, nil
	}
	taskId, err = storestrconv.ParseInt(string(rows[0]["id"]), 10, 64)
	if err != nil {
		return 0, taskExt{}, false, err
	}
	return taskId, taskExt{
		Iid:     string(rows[0]["ext_iid"]),
		Session: string(rows[0]["ext_session"]),
	}, true, nil
}

// listExt returns the extension columns of every task, keyed by task id. One
// query serves the list endpoints and avoids a per-task lookup (N+1).
func listExt() (map[int64]taskExt, error) {
	s := db.NewAutocommitSession()
	defer s.Close()
	rows, err := s.Query(`SELECT id, ext_iid, ext_session FROM tasks WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	exts := make(map[int64]taskExt, len(rows))
	for _, row := range rows {
		taskId, err := storestrconv.ParseInt(string(row["id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		exts[taskId] = taskExt{
			Iid:     string(row["ext_iid"]),
			Session: string(row["ext_session"]),
		}
	}
	return exts, nil
}

// setTaskExtIid persists the VisitNode Id of a task.
func setTaskExtIid(taskId int64, extIid string) error {
	return updateTaskExt(taskId, `UPDATE tasks SET ext_iid = ? WHERE id = ?`, extIid)
}

// setTaskExtSession persists the Session of a task.
func setTaskExtSession(taskId int64, session string) error {
	return updateTaskExt(taskId, `UPDATE tasks SET ext_session = ? WHERE id = ?`, session)
}

// updateTaskExt writes one extension column and asserts exactly one row changed.
func updateTaskExt(taskId int64, query, value string) error {
	s := db.NewAutocommitSession()
	defer s.Close()
	result, err := s.Exec(query, value, taskId)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return storefmt.Errorf("updating ext column of task %d affected %d rows", taskId, affected)
	}
	return nil
}
