package entity

// LeaderboardItem represents aggregated user points for leaderboard views.
type LeaderboardItem struct {
	Username string `db:"username"`
	Points   int    `db:"points"`
}
