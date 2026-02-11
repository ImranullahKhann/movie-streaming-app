package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"context"
	"github.com/ImranullahKhann/movie-streaming-app/server/store"
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserEmail   string
}

func IssueTokens(email string) (*Tokens, error) {
	now := time.Now().UTC()
	t := &Tokens {
		UserEmail: email,
		JTIAcc: uuid.NewString(),
		JTIRef: uuid.NewString(),
		ExpAcc: now.Add(15 * time.Minute),
		ExpRef: now.Add(7 * 25 * time.Hour),
	}

	// HS256 (HMAC with SHA-256) is a symmetric, keyed-hash algorithm used to sign JWT. It uses a single, shared secret for both generating and verifying signatures, making it fast and suitable for monolithic systems where the same entity creates and validates tokens
	// HMAC (Hash-based Message Authentication Code) is a cryptographic mechanism that combines a hash function e.g sha-256 with a secret shared key to simultaneously verify both the data integrity and authenticity of the message. The sender generates a unique MAC (hash) using the message and key; if a receiver's recalculation matches the received MAC, it confirms the message is untampered (integrity) and originated from a trusted source (authenticity)
	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   email,
		ID:        t.JTIAcc,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: email,
		ID: t.JTIRef,
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	t.Refresh, err = ref.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Persist(ctx context.Context, r *store.Redis, t *Tokens) error {
	if err := r.SetJTI(ctx, "access:"+t.JTIAcc, t.UserEmail, t.ExpAcc); err != nil {
		return err
	}
	if err := r.SetJTI(ctx, "refresh:"+t.JTIRef, t.UserEmail, t.ExpRef); err != nil {
		return err
	}
	return nil
}

func SetAuthCookies(c *gin.Context, t *Tokens) {
	// sets a global policy on how the browser should send cookies during cross-site requests
	// SameSiteLaxMode is a modern standard that allows the cookie to be sent when a user clicks a link from an external site to yours, but prevents it from being sent in hidden background requests, like from an img tag or script from another site
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", t.Access, int(time.Until(t.ExpAcc).Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", t.Refresh, int(time.Until(t.ExpRef).Seconds()), "/", "", false, true)
	// options after time
	// "/" : the path, means the cookie is valid for all routes on our domain
	// "" : the domain, leaving this blank means it defaults to the current domain of your api
	// bool (Secure): This ensures that the cookie is only sent over encrypted HTTPS connections, turn false when running locally
	// bool (HttpOnly): This makes the cookie HTTP Only. Meaning it prevents javascript on the client side from accessing the cookie. This makes it harder for an attacker to steal the token via an XSS attack	
}

func ClearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}

func ParseAccess(tokenStr string) (*jwt.RegisteredClaims, error) {
	secret := os.Getenv("ACCESS_SECRET")
	return parseWithSecret(tokenStr, secret)
}


func ParseRefresh(tokenStr string) (*jwt.RegisteredClaims, error) {
	secret := os.Getenv("REFRESH_SECRET")
	return parseWithSecret(tokenStr, secret)
}

func parseWithSecret(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	// telling the parser to only accept HS256
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Extra safety: ensure HMAC family
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
