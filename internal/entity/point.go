package entity

import "time"

type Point struct {
	UserId int       `db:"user_id"`
	TaskId int       `db:"task_id"`
	Points int       `db:"points"`
	UpdAt  time.Time `db:"upd_at"`
}
