package main

import (
	"strings"
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

// YapilacakAdim (Adım)
type YapilacakAdim struct {
	ID               int
	ListeID          int
	Icerik           string
	TamamlandiMi     bool
	OlusTarihi       time.Time
	GuncellemeTarihi time.Time
	SilinmeTarihi    *time.Time
}

// YapilacakListe (Liste)
type YapilacakListe struct {
	ID                int
	Isim              string
	OlusTarihi        time.Time
	GuncellemeTarihi  time.Time
	SilinmeTarihi     *time.Time
	TamamlanmaYuzdesi float64
	Adimlar           []YapilacakAdim
}

// Veri Tabanı zorunlu olmadığı için veriyi RAM'de tutuyorum
var yapilacakListeler = []YapilacakListe{
	{
		ID:               1,
		Isim:             "user1 Alinacaklar Listesi",
		OlusTarihi:       time.Now(),
		GuncellemeTarihi: time.Now(),
		Adimlar: []YapilacakAdim{
			{ID: 1, ListeID: 1, Icerik: "Ekmek al", TamamlandiMi: false, OlusTarihi: time.Now(), GuncellemeTarihi: time.Now()},
			{ID: 2, ListeID: 1, Icerik: "Süt al", TamamlandiMi: true, OlusTarihi: time.Now(), GuncellemeTarihi: time.Now()},
			{ID: 3, ListeID: 1, Icerik: "Şeker al", TamamlandiMi: false, OlusTarihi: time.Now(), GuncellemeTarihi: time.Now()},
		},
	},
	{
		ID:               2,
		Isim:             "admin Yapilacaklar Listesi",
		OlusTarihi:       time.Now(),
		GuncellemeTarihi: time.Now(),
		Adimlar: []YapilacakAdim{
			{ID: 1, ListeID: 2, Icerik: "Projeyi bitir", TamamlandiMi: false, OlusTarihi: time.Now(), GuncellemeTarihi: time.Now()},
			{ID: 2, ListeID: 2, Icerik: "Toplantiya katil", TamamlandiMi: true, OlusTarihi: time.Now(), GuncellemeTarihi: time.Now()},
		},
	},
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

// Token doğrulama
func Dogrulama() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Authorization başlığından tokeni alır
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Token bulunamadi"})
			c.Abort()
			return
		}
		//Doğrulama kısmına tokeni Bearer'den ayırma tokenString değişkenine atar
		hmm := strings.Split(authHeader, " ")
		tokenString := hmm[1]
		//token çözümleme ve doğrulama işlemi
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return Token, nil
		})
		//token geçersizse kullanıcıya hata döndürme
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Geçersiz token"})
			c.Abort()
			return
		}
		//claims ile kullanıcının bilgilerini alma admin mi user1 mi
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Token bilgileri okunamadı"})
			c.Abort()
			return
		}
		//bunu değişkenlere atama
		c.Set("username", claims["sub"])
		c.Set("user_type", claims["user_type"])
		c.Next()
	}
}

// Bütün todosları getirme(kullanıcıya göre)
func GetTodos(c *gin.Context) {
	//Doğrulama fonksiyonundaki kullanıcı bilgilerini değişkene atama
	username := c.GetString("username")
	userType := c.GetString("user_type")
	var aktifListeler []YapilacakListe

	// Eğer admin ise tüm listeyi görsün, user1 ise sadece kendi listelerini görsün
	if userType == "admin" {
		// Admin, tüm listeyi görebilir
		for i := 0; i < len(yapilacakListeler); i++ {
			liste := yapilacakListeler[i]
			if liste.SilinmeTarihi == nil {
				// Tamamlanma yüzdesini hesapla
				var aktifAdimlar []YapilacakAdim
				//Yüzde hesaplama için tamamlanmış adımları alma
				var TamamlanmisAdimlar int
				for j := 0; j < len(liste.Adimlar); j++ {
					adim := liste.Adimlar[j]
					if adim.TamamlandiMi {
						TamamlanmisAdimlar++
					}
					if adim.SilinmeTarihi == nil {
						aktifAdimlar = append(aktifAdimlar, adim)
					}
				}

				// Yüzdeyi hesapla, eğer adım yoksa sıfır kabul et
				if len(liste.Adimlar) > 0 {
					liste.TamamlanmaYuzdesi = float64(TamamlanmisAdimlar) / float64(len(liste.Adimlar)) * 100
				} else {
					liste.TamamlanmaYuzdesi = 0
				}
				liste.Adimlar = aktifAdimlar
				aktifListeler = append(aktifListeler, liste)
			}
		}
	} else {
		// User1 sadece kendi listelerini görebilir
		for i := 0; i < len(yapilacakListeler); i++ {
			liste := yapilacakListeler[i]
			//Liste ismini boşluklara göre ayırma
			listKullanici := strings.Split(liste.Isim, " ")
			if listKullanici[0] == username && liste.SilinmeTarihi == nil {
				// Tamamlanma yüzdesini hesapla
				var aktifAdimlar []YapilacakAdim
				var TamamlanmisAdimlar int
				for j := 0; j < len(liste.Adimlar); j++ {
					adim := liste.Adimlar[j]
					if adim.TamamlandiMi {
						TamamlanmisAdimlar++
					}
					if adim.SilinmeTarihi == nil {
						aktifAdimlar = append(aktifAdimlar, adim)
					}
				}

				// Yüzdeyi hesapla, eğer adım yoksa sıfır kabul et
				if len(liste.Adimlar) > 0 {
					liste.TamamlanmaYuzdesi = float64(TamamlanmisAdimlar) / float64(len(liste.Adimlar)) * 100
				} else {
					liste.TamamlanmaYuzdesi = 0
				}
				liste.Adimlar = aktifAdimlar
				aktifListeler = append(aktifListeler, liste)
			}
		}
	}

	c.JSON(200, gin.H{"kullanici": username, "listeler": aktifListeler})
}

// Liste ekleme
func ListeEkle(c *gin.Context) {
	username := c.GetString("username")
	var listeIsmi string
	var liste YapilacakListe
	if err := c.BindJSON(&listeIsmi); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}
	//Burada listeyi kim oluşturuyorsa listenin başına onun userType'nı ekliyoruz ki liste sıralamada kolaylık sağlasın
	if username == "admin" {
		liste.Isim = "admin " + listeIsmi
	} else if username == "user1" {
		liste.Isim = "user1 " + listeIsmi
	}
	//ID, oluşturma tarihi ve Güncelleme tarihini atama
	liste.ID = len(yapilacakListeler) + 1
	liste.OlusTarihi = time.Now()
	liste.GuncellemeTarihi = time.Now()

	yapilacakListeler = append(yapilacakListeler, liste)

	c.JSON(201, gin.H{"message": "Liste başariyla eklendi", "liste": liste})

}

// liste silme
func ListeSil(c *gin.Context) {
	username := c.GetString("username")
	var ID int
	if err := c.BindJSON(&ID); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}
	zaman := time.Now()
	for i := 0; i < len(yapilacakListeler); i++ {
		liste := yapilacakListeler[i]
		listKullanici := strings.Split(liste.Isim, " ")

		if liste.ID == ID {
			//Kullanıcı admin ise listeyi silsin
			if username == "admin" {
				yapilacakListeler[i].SilinmeTarihi = &zaman
				c.JSON(200, gin.H{"message": "Liste başariyla silindi"})
				return
				//kullanıcı user1 ise ve liste user 1 listesi ise silsin
			} else if username == listKullanici[0] {
				yapilacakListeler[i].SilinmeTarihi = &zaman
				c.JSON(200, gin.H{"message": "Liste başariyla silindi"})
				return
				//kullanıcı user1 ise ve liste admin listesiyse silemesin
			} else {
				c.JSON(403, gin.H{"error": "Bu listeyi silmeye yetkiniz yok"})
				return
			}
		}
	}
	c.JSON(404, gin.H{"error": "Liste bulunamadı"})
}
func main() {
	r := gin.Default()
	r.POST("/login", Giris)

	// Korumalı alan
	protected := r.Group("/")
	protected.Use(Dogrulama())
	protected.GET("/todos", GetTodos)
	protected.POST("/lists", ListeEkle)
	protected.DELETE("/lists", ListeSil)
	r.Run()
}
