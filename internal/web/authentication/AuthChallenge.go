package authentication

import "time"

type AuthChallenge struct {
	Expiration time.Time
}
