package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID, _ := uuid.NewRandom()
	tokenSecret := "secret_token"
	dur := 5 * time.Hour

	jwtString, err := MakeJWT(userID, tokenSecret, dur)
	if err != nil {
		t.Errorf("[TEST ERROR] MakeJWT: %v\n", err)
	}

	validatedUserID, err := ValidateJWT(jwtString, tokenSecret)
	if err != nil {
		t.Errorf("[TEST ERROR] ValidateJWT: %v\n", err)
	}

	if userID != validatedUserID {

		t.Errorf("%s != %s user ids are not equal", userID, validatedUserID)
	}

}
