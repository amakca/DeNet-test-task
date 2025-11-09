package entity

import "time"

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	Referrer  *string   `db:"referrer"`
	Email     *string   `db:"email"`
}
