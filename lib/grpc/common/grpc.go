package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type challengeClaims struct {
	ChallengeID string `json:"id"`
	Time        int64  `json:"t,omitempty"`
	Static      bool   `json:"s,omitempty"`
}

// The implementation is not relevant here since we don't use standard claims.

func (c challengeClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(0, 0)), nil
}

func (c challengeClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(0, 0)), nil
}

func (c challengeClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c challengeClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c challengeClaims) GetSubject() (string, error) {
	return "", nil
}

func (c challengeClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func ToStatus(status string) Status {
	if s, ok := Status_value[status]; ok {
		return Status(s)
	}
	return Status_Unspecified
}

func (c *ChallengeResponse) CalculateAnimatedToken(start time.Time) (string, error) {
	diff := time.Since(start)
	seconds := int64(diff.Seconds())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &challengeClaims{
		ChallengeID: c.ChallengeId,
		Time:        seconds,
		Static:      false,
	})
	t, err := token.SignedString([]byte(c.ChallengeSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (c *ChallengeResponse) CalculateStaticToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &challengeClaims{
		ChallengeID: c.ChallengeId,
		Static:      true,
	})
	t, err := token.SignedString([]byte(c.ChallengeSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}
