package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenPayload struct {
	UserId    int
	CreatedAt int
}

const secret = "secret"
const month_duration = 30 * 24 * 60 * 60 * 1000

var ErrInvalidToken = errors.New("invalid token")

func encode(userId int) string {
	hash := sha256.New()
	createdAt := int(time.Now().UnixMilli())
	payloadStr := fmt.Sprintf("%d.%d", userId, createdAt)

	hash.Write([]byte(payloadStr))

	controlSign := hash.Sum([]byte(secret))
	encodedSign := base64.StdEncoding.EncodeToString(controlSign)

	return fmt.Sprintf("%s.%s", payloadStr, encodedSign)
}

func decode(token string) (tokenPayload, error) {
	tp := tokenPayload{}
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return tp, ErrInvalidToken
	}

	userId, err := strconv.Atoi(parts[0])
	if err != nil {
		return tp, ErrInvalidToken
	}

	createdAt, err := strconv.Atoi(parts[1])
	if err != nil {
		return tp, ErrInvalidToken
	}

	diff := math.Abs(float64(int(time.Now().UnixMilli()) - createdAt))
	if diff > month_duration {
		return tp, ErrInvalidToken
	}

	sign, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return tp, ErrInvalidToken
	}

	hash := sha256.New()
	payloadStr := fmt.Sprintf("%d.%d", userId, createdAt)

	hash.Write([]byte(payloadStr))

	controlSign := hash.Sum([]byte(secret))
	if equal := bytes.Equal(sign, controlSign); !equal {
		return tp, ErrInvalidToken
	}

	tp.UserId = userId
	tp.CreatedAt = createdAt

	return tp, nil
}

func CreateToken(userId int) string {
	return encode(userId)
}

func ExtractUserId(ctx *gin.Context) int {
	tp, exists := ctx.Get(AuthCookieName)
	if !exists {
		return -1
	}

	payload, ok := tp.(tokenPayload)
	if !ok {
		return -1
	}

	return payload.UserId
}
