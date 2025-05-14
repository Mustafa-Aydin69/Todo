# TO-DO App

Bu proje, kullanıcıların görevlerini takip edebilmeleri için bir **REST API** sağlar. Kullanıcılar, görev listelerini oluşturabilir, silebilir ve güncelleyebilir. Aynı şekilde her listeye adımlar ekleyebilir, silebilir veya güncelleyebilirler. API, güvenlik amacıyla **JSON Web Token (JWT)** kullanır, böylece yalnızca geçerli bir token ile korunmuş işlemler yapılabilir.

## Özellikler
## 🌐 Canlı Yayın Linki
[https://todoapp-3s7x.onrender.com](https://todoapp-3s7x.onrender.com)
- **Kullanıcı Girişi**: Kullanıcı adı ve şifre ile giriş yapılır. Başarılı giriş sonrası bir **token** verilir.
- **Görev Listesi Yönetimi**: Görev listeleri oluşturulabilir, güncellenebilir veya silinebilir.
- **Adım Yönetimi**: Her görev listesine adımlar eklenebilir, güncellenebilir veya silinebilir.
- **Koruma**: Görev listesi ve adım yönetim işlemleri **token** ile korunur. Yalnızca geçerli bir token ile erişilebilir. admin bütün Todolara erişebilirken user1 sadece kendi Todolarına erişebilir.

## Gereksinimler

- **GoLang**: API'yi geliştirmek için Go dilini kullandım.
- **Thunder Client** veya **Postman**: API'yi test etmek için **Thunder Client** eklentisi kullanılabilir. (alternatif olarak Postman kullanılabilir).
- **JSON Web Token (JWT)**: Token doğrulaması için kullanılmıştır.

## API İşlemleri ve Açıklamaları
- **Kayıtlı Kullanıcılar**: Sistemde Kayıtlı olan kullanıcılar user1 ve admindir.
- ![Kayıtlı Kullanıcılar](https://github.com/user-attachments/assets/ed378ee8-2924-458f-8263-778d72d216dc)
- **Giriş İşlemi ve Token Alma**: **POST** metodu, sunucuya veri göndermek için kullanılır. İstek, `/login` endpoint'ine yapılmıştır. Kullanıcı adı ve şifre bilgileri **JSON** formatında sunucuya gönderilmiştir:
  ![Image](https://github.com/user-attachments/assets/5133df68-dcf3-4756-84c9-6d7132e2e788)
- **Get Todos**: **GET** metodu, sunucudan veri almak için kullanılır. `/todos` endpoint'ine bir GET isteği gönderilmektedir. Bu endpoint, giriş yapan kullanıcının yapılacaklar listesine ait verileri almak için kullanılır.Bir **Authorization** başlığı bulunmaktadır. Bu başlık, API'ye erişim için kimlik doğrulama token'ını içerir. "Bearer" kelimesinin ardından gelen uzun karakter dizisi ise geçerli token'dır.
  ![Image](https://github.com/user-attachments/assets/85d25b86-3254-4003-ac63-0038cf8550a8)
- **Liste Ekle**: **POST** metodu sunucuya yeni bir veri eklemek için kullanılan bir metodtur. Burada `/lists` endpoint'ine yapılmış bir istek var, bu yeni bir liste oluşturma işlemi içindir.  **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Liste ismi  **JSON** formatında sunucuya gönderilir.
![Image](https://github.com/user-attachments/assets/b63afbad-edd5-4f57-a22a-7916e69011f7) ![Image](https://github.com/user-attachments/assets/4281b9c0-aba9-4e39-b629-934e3942559f)
- **Liste Sil**: **DELETE** metodu sunucuda bulunan belirli bir kaynağı silmek için kullanılır. Burada `/lists` endpoint'ine bir DELETE isteği yapılır. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Silinecek listenin ID'si  **JSON** formatında sunucuya gönderilmiştir.
![Image](https://github.com/user-attachments/assets/f05221d5-dcf4-4515-a81f-44eb258dd0af)  ![Image](https://github.com/user-attachments/assets/2c1b46ff-fe40-4fd4-b02a-479939c10488)  
- **Liste Güncelleme**: **PUT** metodu sunucudaki mevcut bir kaynağı güncellemek için kullanılır. Burada `/lists` endpoint'ine bir PUT isteği yapılmış ve var olan bir listenin bilgileri güncellenmek isteniyor. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Gönderilen **JSON** verisinde, güncellenmek istenen listenin **ID**'sini alır ve **İsim** alanını günceller.
![Image](https://github.com/user-attachments/assets/d72e0069-3b3b-48d3-9eed-4748f93dee56) ![Image](https://github.com/user-attachments/assets/b70ef4bf-d1b8-4342-be57-48c34bf45bc9)
- **Adım Ekle**: **POST** metodu sunucuya yeni bir veri eklemek için kullanılır. Burada `/steps` endpoint'ine yapılan istek, yeni bir adım eklemek amacıyla gönderilmiştir. Gönderilen **JSON** verisi, eklenmek istenen adımı temsil eder. Alınan Liste ID'sine İçerik adımı eklenir.
 ![Image](https://github.com/user-attachments/assets/cdb5cff7-7a27-488b-8ed4-f2a3ab7797b9) ![Image](https://github.com/user-attachments/assets/0e6d189d-1727-4e8e-80f3-4c39ee185899)
- **Adım Sil**: **DELETE** metodu sunucuda bulunan belirli bir kaynağı silmek için kullanılır. Burada `/steps` endpoint'ine bir Delete isteği yapılır. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Silinecek Adımın ve hangi listede olduğu bilgileri **JSON** formatında sunucuya gönderilir.
  ![Image](https://github.com/user-attachments/assets/4b2f0b8c-c6be-4103-912e-df9b4dcd0334) ![Image](https://github.com/user-attachments/assets/f21b00e9-a956-4dff-90fb-39f1f2d4a28c)
- **Adım Güncelle**: **PUT** metodu sunucudaki mevcut bir adımı güncellemek için kullanılır.Burada `/steps` endpoint'ine bir PUT isteği yapılır ve var olan bir adımın bilgileri güncellenir. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Gönderilen **JSON** verisinde, güncellenmek istenen adımın ListeID'si ve ID'si üzerinden İçeriği ve Tamamlanma durumu güncellenir.
  ![Image](https://github.com/user-attachments/assets/4f365fc8-be45-425a-8142-5efcff2f8e05) ![Image](https://github.com/user-attachments/assets/954d14d7-01be-431b-b411-b696a9e09b2c)
