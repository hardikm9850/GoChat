package config

import (
    "os"
    "time"

    "github.com/hardikm9850/authkit/jwt"
)

func JWT() jwt.Config {
    return jwt.Config{
        Algorithm:      jwt.HS256,
        Secret:         os.Getenv("JWT_SECRET"),
        Issuer:         os.Getenv("JWT_ISSUER"),
        Audience:       os.Getenv("JWT_AUDIENCE"),
        AccessTokenTTL: 24 * time.Hour,
    }
}
