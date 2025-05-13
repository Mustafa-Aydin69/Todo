package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var Token = []byte("Kullanici_Tokeni")

// Örnek kullanıcılar
var Kullanicilar = map[string]string{
	"user1": "1234",
	"admin": "admin",
}

func TokenOlustur(username, userType string) (string, error) {
	claims := &MyClaims{
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Token)
}

// Giris işlemi
func Giris(c *gin.Context) {
	var Kontrol struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&Kontrol); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}

	Sifre, Varmi := Kullanicilar[Kontrol.Username]
	if !Varmi || Sifre != Kontrol.Password {
		c.JSON(401, gin.H{"error": "Geçersiz kullanici adi veya şifre"})
		return
	}

	// Kullanıcı tipi (admin veya user1) belirlenmeli
	var KullaniciTipi string
	if Kontrol.Username == "admin" {
		KullaniciTipi = "admin"
	} else {
		KullaniciTipi = "user1"
	}

	// Token üretimi
	token, err := TokenOlustur(Kontrol.Username, KullaniciTipi)
	if err != nil {
		c.JSON(500, gin.H{"error": "Token üretilemedi"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

type MyClaims struct {
	Username string `json:"sub"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func main() {
	r := gin.Default()
	r.POST("/login", Giris)

	r.Run()
}
