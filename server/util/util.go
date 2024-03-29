package util

import (
	// import go jwt library
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zamachnoi/viewthis/data"
	"github.com/zamachnoi/viewthis/models"
)

type ContextKey string

const UserKey ContextKey = "user"

const (
	// discord refresh token url
	DiscordRefreshTokenURL = "https://discord.com/api/v10/oauth2/token"
)
type SessionJWT struct {
    DiscordID string `json:"discord_id"`
    Avatar    string `json:"avatar"`
    Username  string `json:"username"`
    DBID      uint   `json:"db_id"`
    Premium   bool   `json:"premium"`
}

// Create Struct to get Subject from the token
type SessionJWTWithClaims struct {
    SessionJWT
	jwt.RegisteredClaims
}

func ParseJWTClaims(tokenString string) (*jwt.Token, SessionJWTWithClaims, error) {
    claims := SessionJWTWithClaims{}
    token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return token, claims, err
}

func GenerateSessionJWT(sessionJwt SessionJWT) (string, error) {
	claims := SessionJWTWithClaims{
		SessionJWT: sessionJwt,
        RegisteredClaims: jwt.RegisteredClaims{
            NotBefore: jwt.NewNumericDate(time.Now()),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(GetJWTExpiry()),
        },
	}

	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }
	return tokenString, err
}

func GetJWTExpiry() time.Time {
	return time.Now().Add(time.Hour * 2) // set to 1 day for testing purposes :/ TODO CHANGE LATER
}

func GetCookieExpiry() time.Time {
    return time.Now().Add(time.Hour * (24 * 14)) // 2 weeks TODO: make env var
}

func EncryptRefreshToken(token string) (string, error) {
    // AES ENCRYPT CTR
    block, err := aes.NewCipher([]byte(os.Getenv("AES_ENCRYPTION_KEY")))
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(token))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(token))

    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptRefreshToken(encryptedToken string) (string, error) {
    // AES DECRYPT
    block, err := aes.NewCipher([]byte(os.Getenv("AES_ENCRYPTION_KEY")))
    if err != nil {
        return "", err
    }

    decodedToken, err := base64.URLEncoding.DecodeString(encryptedToken)
    if err != nil {
        return "", err
    }

    if len(decodedToken) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv, ciphertext := decodedToken[:aes.BlockSize], decodedToken[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}

func EncodeDiscordUserInfo(discordUser models.DiscordUser, refreshToken string) (*models.User, error) {
    // store discord info in db
    newUserInfo, err := data.GetUserByDiscordID(discordUser.ID)
    if err != nil {
        return nil, err
    }

    newUserInfo.Username = discordUser.Username
    newUserInfo.Avatar = discordUser.Avatar
    newUserInfo.DiscordID = discordUser.ID
    newUserInfo.RefreshExpiry = time.Now().Add(time.Hour * 24 * 30)
    newUserInfo.RefreshToken, err = EncryptRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }

    return newUserInfo, nil
}

func SetJWTCookie(jwt string, w http.ResponseWriter) {
	expiry := GetCookieExpiry()
    ExpireCookie("_viewthis_jwt", w)

    domain := ".viewthis.app"
    if os.Getenv("DEV") == "true" {
        domain = "localhost"
    }

    // if the domain is localhost, make secure false
    secure := true
    if os.Getenv("DEV") == "true" {
        secure = false
    }
    
    http.SetCookie(w, &http.Cookie{
        Name:     "_viewthis_jwt",
        Value:    jwt,
        Expires:  expiry,
		HttpOnly: true,
        Domain:  domain,
		SameSite: http.SameSiteLaxMode,
		Secure: secure,
        Path:    "/", // set to root so it's accessible from all pages
    })
}

func ExpireCookie(name string, w http.ResponseWriter) {
 
    // Create a cookie with the same name, path, and domain
    cookie := http.Cookie{
        Name:   name,
        Path:   "/", // Ensure this matches the path of the cookie you're clearing
        Expires: time.Unix(0, 0), // Set the expiration date to the past
        MaxAge:   -1, // Tells the browser to delete the cookie immediately
        HttpOnly: true,
    }
    http.SetCookie(w, &cookie)
}

// Helper function to parse limit and page parameters
func ParseLimitAndPage(r *http.Request) (int, int) {
    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil || limit <= 0 {
        limit = 9 // default limit
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page <= 0 {
        page = 1 // default page
    }

    return limit, page
}

// Helper function to get JWT value
func GetJWTValue(r *http.Request) string {
    jwtCookie, err := r.Cookie("_viewthis_jwt")
    log.Println("Getting JWT value: ", jwtCookie)
    if err != nil {
        return ""
    }
    return jwtCookie.Value
}