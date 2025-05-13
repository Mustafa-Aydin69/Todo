# TO-DO App

Bu proje, kullanıcıların görevlerini takip edebilmeleri için bir **REST API** sağlar. Kullanıcılar, görev listelerini oluşturabilir, silebilir ve güncelleyebilir. Aynı şekilde her listeye adımlar ekleyebilir, silebilir veya güncelleyebilirler. API, güvenlik amacıyla **JSON Web Token (JWT)** kullanır, böylece yalnızca geçerli bir token ile korunmuş işlemler yapılabilir.

## Özellikler

- **Kullanıcı Girişi**: Kullanıcı adı ve şifre ile giriş yapılır. Başarılı giriş sonrası bir **token** verilir.
- **Görev Listesi Yönetimi**: Görev listeleri oluşturulabilir, güncellenebilir veya silinebilir.
- **Adım Yönetimi**: Her görev listesine adımlar eklenebilir, güncellenebilir veya silinebilir.
- **Koruma**: Görev listesi ve adım yönetim işlemleri **token** ile korunur. Yalnızca geçerli bir token ile erişilebilir. admin bütün Todolara erişebilirken user1 sadece kendi Todolarına erişebilir.

## Gereksinimler

- **GoLang**: API'yi geliştirmek için Go dilini kullandım.
- **Thunder Client** veya **Postman**: API'yi test etmek için **Thunder Client** eklentisi kullandım. (alternatif olarak Postman kullanılabilir).
- **JSON Web Token (JWT)**: Token doğrulaması için kullanılmıştır.

## API İşlemleri ve Açıklamaları
- **Kayıtlı Kullanıcılar**: Sistemde Kayıtlı olan kullanıcılar user1 ve admindir.
- ![Kayıtlı Kullanıcılar](https://github.com/user-attachments/assets/ed378ee8-2924-458f-8263-778d72d216dc)
- **Giriş İşlemi ve Token Alma**: **POST** metodu, sunucuya veri göndermek için kullanılır. İstek, `/login` endpoint'ine yapılmıştır. Kullanıcı adı ve şifre bilgileri **JSON** formatında sunucuya gönderilmiştir:
  ![Giriş ve Token](https://github.com/user-attachments/assets/d2f14298-f9b1-47af-9b06-04bff1ba3554)
- **Get Todos**: **GET** metodu, sunucudan veri almak için kullanılır. `/todos` endpoint'ine bir GET isteği gönderilmektedir. Bu endpoint, giriş yapan kullanıcının yapılacaklar listesine ait verileri almak için kullanılır.Bir **Authorization** başlığı bulunmaktadır. Bu başlık, API'ye erişim için kimlik doğrulama token'ını içerir. "Bearer" kelimesinin ardından gelen uzun karakter dizisi ise geçerli token'dır.
  ![GetTodos](https://github.com/user-attachments/assets/5546335b-4378-4417-849c-91c72b847244)
- **Liste Ekle**: **POST** metodu sunucuya yeni bir veri eklemek için kullanılan bir metodtur. Burada `/lists` endpoint'ine yapılmış bir istek var, bu yeni bir liste oluşturma işlemi içindir.  **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Liste ismi  **JSON** formatında sunucuya gönderilir.
![Image](https://github.com/user-attachments/assets/0b5e3979-001b-4a8f-af58-9c8ba791970f) ![Image](https://github.com/user-attachments/assets/a2aa3b31-502f-411a-b548-8e2fa82a9649)
- **Liste Sil**: **DELETE** metodu sunucuda bulunan belirli bir kaynağı silmek için kullanılır. Burada `/lists` endpoint'ine bir DELETE isteği yapılır. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Silinecek listenin ID'si  **JSON** formatında sunucuya gönderilmiştir.
![Image](https://github.com/user-attachments/assets/b9ad1e3d-ca5d-4d01-889a-069d5e81c4d5)  ![Image](https://github.com/user-attachments/assets/49463609-8f07-484f-b56d-58769dde344b)  
- **Liste Güncelleme**: **PUT** metodu sunucudaki mevcut bir kaynağı güncellemek için kullanılır. Burada `/lists` endpoint'ine bir PUT isteği yapılmış ve var olan bir listenin bilgileri güncellenmek isteniyor. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Gönderilen **JSON** verisinde, güncellenmek istenen listenin **ID**'sini alır ve **İsim** alanını günceller.
![Image](https://github.com/user-attachments/assets/2a7a5e49-7da9-4a75-aab9-59e77162e5ff) ![Image](https://github.com/user-attachments/assets/0a6fb169-9c5c-4d61-be18-9c39df27e551)
- **Adım Ekle**: **POST** metodu sunucuya yeni bir veri eklemek için kullanılır. Burada `/steps` endpoint'ine yapılan istek, yeni bir adım eklemek amacıyla gönderilmiştir. Gönderilen **JSON** verisi, eklenmek istenen adımı temsil eder. Alınan Liste ID'sine İçerik adımı eklenir.
 ![Image](https://github.com/user-attachments/assets/68d0c9a2-d9ed-4199-91e0-6fe254183b5a) ![Image](https://github.com/user-attachments/assets/85f9d4a8-2a09-42ba-9274-233aef1157f3)
- **Adım Sil**: **DELETE** metodu sunucuda bulunan belirli bir kaynağı silmek için kullanılır. Burada `/steps` endpoint'ine bir Delete isteği yapılır. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Silinecek Adımın ve hangi listede olduğu bilgileri **JSON** formatında sunucuya gönderilir.
  ![Image](https://github.com/user-attachments/assets/0f2943d1-44e0-4781-8b57-1cce7b1d20e1) ![Image](https://github.com/user-attachments/assets/29c418fa-aeca-4aa2-855f-ef4b755b2ab0)
- **Adım Güncelle**: **PUT** metodu sunucudaki mevcut bir adımı güncellemek için kullanılır.Burada `/steps` endpoint'ine bir PUT isteği yapılır ve var olan bir adımın bilgileri güncellenir. **Authorization** başlığı kullanılarak kimlik doğrulama yapılır. Gönderilen **JSON** verisinde, güncellenmek istenen adımın ListeID'si ve ID'si üzerinden İçeriği ve Tamamlanma durumu güncellenir.
  ![Image](https://github.com/user-attachments/assets/97acf725-1760-4326-a0b9-52eb3a944c6e) ![Image](https://github.com/user-attachments/assets/c4330aa0-9336-4785-9d04-aeec5046d033)
