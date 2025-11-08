package entity

type Task struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Descr string `db:"descr"`
}

const (
	TaskSubscribeToTelegram = "subscribe_to_telegram"
	TaskSubscribeToTwitter  = "subscribe_to_twitter"
	TaskInviteFriend        = "invite_friend"
	TaskEnterReferralCode   = "enter_referral_code"
	TaskEnterEmail          = "enter_email"
)
