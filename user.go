package apple

import (
	"github.com/tideland/gorest/jwt"
)

var (
	// RealUserStatusUnsupported unsupported, only works in iOS >= 14.
	RealUserStatusUnsupported RealUserStatus = 0
	// RealUserStatusUnknown cannot determine if the user is real.
	RealUserStatusUnknown RealUserStatus = 1
	// RealUserStatusLikelyReal user is likely real.
	RealUserStatusLikelyReal RealUserStatus = 2
)

// RealUserStatus an integer value that indicates whether the user appears to be
// a real person.
type RealUserStatus int

// AppleUser is the model to hold information about the user.
type AppleUser struct {
	// UID Apple unique identification for the user.
	UID string `json:"uid"`

	// TransferUID Apple unique transfer identification for the user.
	TransferUID string `json:"transfer_uid"`

	// Email Apple user email.
	Email string `json:"email"`

	// EmailVerified whether the email is verified.
	EmailVerified bool `json:"email_verified"`

	// IsPrivateEmail whether the email shared by the user is the proxy address.
	IsPrivateEmail bool `json:"is_private_email"`

	// RealUserStatus an integer value that indicates whether the user appears
	// to be a real person.
	RealUserStatus RealUserStatus `json:"real_user_status"`
}

// GetUserInfoFromIDToken retrieve the user info from the JWT id token.
func GetUserInfoFromIDToken(idToken string) (*AppleUser, error) {
	token, err := jwt.Decode(idToken)
	if err != nil {
		return nil, err
	}

	u := AppleUser{}
	claims := token.Claims()
	if sub, ok := claims["sub"].(string); ok {
		u.UID = sub
	}

	if transferSub, ok := claims["transfer_sub"].(string); ok {
		u.TransferUID = transferSub
	}

	if email, ok := claims["email"].(string); ok {
		u.Email = email
	}

	if emailVerified, ok := claims["email_verified"].(bool); ok {
		u.EmailVerified = emailVerified
	}

	if isPrivateEmail, ok := claims["is_private_email"].(bool); ok {
		u.IsPrivateEmail = isPrivateEmail
	}

	// As in JSON ints and floats are the same type, the number type, we must check if the number is
	// either an int or a float, and convert it to the first if the later.
	var (
		realUserStatus int
		ok             bool
	)
	if realUserStatus, ok = claims["real_user_status"].(int); !ok {
		var realUserStatusFloat float64
		if realUserStatusFloat, ok = claims["real_user_status"].(float64); ok {
			realUserStatus = int(realUserStatusFloat)
		}
	}

	if ok {
		switch realUserStatus {
		case int(RealUserStatusLikelyReal):
			u.RealUserStatus = RealUserStatusLikelyReal
		case int(RealUserStatusUnknown):
			u.RealUserStatus = RealUserStatusUnknown
		default:
			u.RealUserStatus = RealUserStatusUnsupported
		}
	}

	return &u, nil
}
