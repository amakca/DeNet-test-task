package entity

import "time"

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	Points    int       `db:"points"`
	Referrer  string    `db:"referrer"`
}

type Task struct {
	Id          int       `db:"id"`
	Type        string    `db:"type"`
	CompletedAt time.Time `db:"completed_at"`
	UserId      int       `db:"user_id"`
}
