package entity

type Task struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Descr  string `db:"descr"`
	Points int    `db:"points"`
}
