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

// liste isim güncelleme
func ListeGünc(c *gin.Context) {
	username := c.GetString("username")
	//JSON üzerinden değiştilecek listenin ID'sini ve Yeni ismi al
	var data struct {
		ID   int    `json:"ID"`
		Isim string `json:"Isim"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}

	for i := 0; i < len(yapilacakListeler); i++ {
		liste := yapilacakListeler[i]
		listKullanici := strings.Split(liste.Isim, " ")
		if liste.ID == data.ID {
			//Kullanıcı ne ise ona göre güncelleme yapma
			if username == "admin" || username == listKullanici[0] {
				if username == "admin" {
					liste.Isim = "admin" + " " + data.Isim
				} else if username == "user1" {
					liste.Isim = "user1" + " " + data.Isim
				}
				liste.GuncellemeTarihi = time.Now()
				yapilacakListeler[i] = liste
				c.JSON(200, gin.H{"message": "Liste başarıyla güncellendi", "Liste": gin.H{"ID": yapilacakListeler[i].ID, "Isim": yapilacakListeler[i].Isim}})
				return
				//Kullanıcı user1 ise ve Liste IDsi admin listesiyse hata döndürsün
			} else {
				c.JSON(403, gin.H{"error": "Bu Listeyi güncellemeye yetkiniz yok"})
				return
			}
		}
	}
	c.JSON(404, gin.H{"error": "Liste bulunamadi"})

}

// Listeye adım ekleme
func AdimEkle(c *gin.Context) {
	username := c.GetString("username")
	zaman := time.Now()
	//JSON üzerinden değiştilecek adım eklenecek listenin ID'sini ve içeriğini al
	var adimekleme struct {
		ListeID int    `json:"ListeID"`
		Icerik  string `json:"Icerik"`
	}
	if err := c.BindJSON(&adimekleme); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}
	for i := 0; i < len(yapilacakListeler); i++ {
		liste := yapilacakListeler[i]
		listKullanici := strings.Split(liste.Isim, " ")
		if liste.ID == adimekleme.ListeID {
			//Kullanıcı Admin ise liste ID'si ne ise ona göre adım ekleme yapma
			if username == "admin" {
				yeniAdim := YapilacakAdim{
					ID:               len(liste.Adimlar) + 1, // Yeni adım ID'si
					Icerik:           adimekleme.Icerik,
					ListeID:          adimekleme.ListeID,
					TamamlandiMi:     false,
					OlusTarihi:       zaman,
					GuncellemeTarihi: zaman,
				}
				liste.Adimlar = append(liste.Adimlar, yeniAdim)
				yapilacakListeler[i] = liste
				c.JSON(201, gin.H{"message": "Adım başarıyla eklendi", "adim": yeniAdim})
				return
				//Kullanıcı user1 ise ve Listede user1 listesiyse ekleme yapma
			} else if username == listKullanici[0] {
				yeniAdim := YapilacakAdim{
					ID:               len(liste.Adimlar) + 1, // Yeni adım ID'si
					Icerik:           adimekleme.Icerik,
					ListeID:          adimekleme.ListeID,
					TamamlandiMi:     false,
					OlusTarihi:       zaman,
					GuncellemeTarihi: zaman,
				}
				liste.Adimlar = append(liste.Adimlar, yeniAdim)
				yapilacakListeler[i] = liste
				c.JSON(201, gin.H{"message": "Adım başarıyla eklendi", "adim": yeniAdim})
				return
				//Kullanıcı user1 ve Liste admin listesiyse hata döndürme
			} else {
				c.JSON(403, gin.H{"error": "Bu listeye adım eklemeye yetkiniz yok"})
				return
			}

		}
	}
	c.JSON(404, gin.H{"error": "Liste bulunamadi"})
}

// Listeden adim silme
func AdimSil(c *gin.Context) {
	username := c.GetString("username")
	zaman := time.Now()
	//Silinecek adımın Idsi ve Listenin Idsini alma
	var data struct {
		ID      int `json:"ID"`
		ListeID int `json:"ListeID"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}

	for i := 0; i < len(yapilacakListeler); i++ {
		liste := yapilacakListeler[i]
		listKullanici := strings.Split(liste.Isim, " ")
		if liste.ID == data.ListeID {
			for j := 0; j < len(liste.Adimlar); j++ {
				adim := liste.Adimlar[j]
				if adim.ID == data.ID {
					//Adminse direkt silme, User1 ise Listenin User1 listesi olduğu kontrol edip silme
					if username == "admin" || username == listKullanici[0] {
						liste.Adimlar[j].SilinmeTarihi = &zaman
						yapilacakListeler[i] = liste
						c.JSON(200, gin.H{"message": "Adim başariyla silindi"})
						return
						//Kullanıcı user1 ve Liste admin listesiyse hata döndürme
					} else {
						c.JSON(403, gin.H{"error": "Bu adimi silmeye yetkiniz yok"})
						return
					}
				}
			}
			c.JSON(404, gin.H{"error": "Adim bulunamadi"})
			return
		}
	}
	c.JSON(404, gin.H{"error": "Liste bulunamadi"})
}

// Listedeki Adımi güncelleme
func AdimGünc(c *gin.Context) {
	username := c.GetString("username")
	//Güncellenecek değerleri alma
	var data struct {
		ID           int    `json:"ID"`
		ListeID      int    `json:"ListeID"`
		Icerik       string `json:"Icerik"`
		TamamlandiMi bool   `json:"TamamlandiMi"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}
	for i := 0; i < len(yapilacakListeler); i++ {
		liste := yapilacakListeler[i]
		listKullanici := strings.Split(liste.Isim, " ")
		if liste.ID == data.ListeID {
			for j := 0; j < len(liste.Adimlar); j++ {
				adim := liste.Adimlar[j]
				if adim.ID == data.ID {
					//Adminse direkt güncelleme, User1 ise Listenin User1 listesi olup olmadığını kontrol etme
					if username == "admin" || username == listKullanici[0] {
						liste.Adimlar[j].Icerik = data.Icerik
						liste.Adimlar[j].TamamlandiMi = data.TamamlandiMi
						liste.Adimlar[j].GuncellemeTarihi = time.Now()

						yapilacakListeler[i] = liste

						c.JSON(200, gin.H{"message": "Adım başarıyla güncellendi", "adim": liste.Adimlar[j]})
						return
						//Kullanıcı user1 ve Liste admin listesiyse hata döndürme
					} else {
						c.JSON(403, gin.H{"error": "Bu adimi güncellemeye yetkiniz yok"})
						return
					}

				}
			}
			c.JSON(404, gin.H{"error": "Adim bulunamadi"})
			return
		}
	}
	c.JSON(404, gin.H{"error": "Liste bulunamadi"})
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
	protected.PUT("/lists", ListeGünc)
	protected.POST("/steps", AdimEkle)
	protected.DELETE("/steps", AdimSil)
	protected.PUT("/steps", AdimGünc)
	r.Run()
}
