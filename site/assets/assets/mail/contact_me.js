document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("contactForm");
    const phoneCable = document.querySelector(".phone-cable");
    const smartphoneMockup = document.querySelector(".smartphone-mockup");
    let animationInterval;

    // Kablo takılıp çıkarma animasyonunu başlat
    function startCableAnimation() {
        let isPlugged = false;
        animationInterval = setInterval(() => {
            if (isPlugged) {
                // Çıkar
                phoneCable.classList.remove("plugged");
                smartphoneMockup.classList.remove("charging");
            } else {
                // Tak
                phoneCable.classList.add("plugged");
                smartphoneMockup.classList.add("charging");
            }
            isPlugged = !isPlugged;
        }, 3000); // Her 3 saniyede bir takıp çıkar
    }

    // Sayfa yüklendiğinde animasyonu başlat
    startCableAnimation();


    form.addEventListener("submit", function (e) {
        e.preventDefault();

        // Hata mesajlarını temizle
        const helpBlocks = form.querySelectorAll(".help-block.text-danger");
        helpBlocks.forEach(block => (block.textContent = ""));

        // Form verilerini al
        const name = document.getElementById("name").value.trim();
        const email = document.getElementById("email").value.trim();
        const phone = document.getElementById("phone").value.trim();
        const message = document.getElementById("message").value.trim();

        let hasError = false;

        // Zorunlu alan kontrolü
        if (name === "") {
            document.querySelector("#name + .help-block").textContent = "Lütfen adınızı ve soyadınızı girin.";
            hasError = true;
        }

        if (email === "") {
            document.querySelector("#email + .help-block").textContent = "Lütfen bir e-posta adresi girin.";
            hasError = true;
        } else if (!isValidEmail(email)) {
            document.querySelector("#email + .help-block").textContent = "Geçerli bir e-posta adresi girin.";
            hasError = true;
        }

        if (phone === "") {
            document.querySelector("#phone + .help-block").textContent = "Lütfen telefon numaranızı girin.";
            hasError = true;
        } else if (!isValidPhone(phone)) { // Güncellenmiş fonksiyonu kullanıyoruz
            document.querySelector("#phone + .help-block").textContent = "Geçerli bir telefon numarası girin (sadece rakam veya 555-555-55-55 formatı).";
            hasError = true;
        }

        if (message === "") {
            document.querySelector("#message + .help-block").textContent = "Lütfen mesajınızı girin.";
            hasError = true;
        }

        if (hasError) {
            return;
        }

        // AJAX ile formu gönder
        const data = {
            name: name,
            email: email,
            phone: phone,
            message: message
        };

        const submitButton = document.getElementById("sendMessageButton");
        submitButton.disabled = true;

        fetch("/contact/submit", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(data).toString()
        })
            .then(response => {
                if (response.ok) {
                    // Başarı durumunda
                    document.getElementById('success').innerHTML = "<div class='alert alert-success'><strong>Mesajınız başarıyla gönderildi.</strong></div>";
                    form.reset();
                } else {
                    // Hata durumunda
                    return response.text().then(text => {
                        throw new Error(text || "Bilinmeyen bir hata oluştu.");
                    });
                }
            })
            .catch(error => {
                let errorMessage = "Üzgünüm, mesajınız gönderilemedi. Lütfen daha sonra tekrar deneyin.";
                if (error.message.includes("zorunlu alanları")) {
                    errorMessage = "Lütfen tüm zorunlu alanları doldurun.";
                } else if (error.message.includes("geçerli bir e-posta")) {
                    errorMessage = "Geçerli bir e-posta adresi girin.";
                } else if (error.message.includes("Telefon numarası formatı")) {
                    errorMessage = "Telefon numarası formatı geçerli değil (örn: 555-555-55-55).";
                }
                document.getElementById('success').innerHTML = "<div class='alert alert-danger'><strong>" + errorMessage + "</strong></div>";
            })
            .finally(() => {
                submitButton.disabled = false;
            });
    });

    function isValidEmail(email) {
        const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return re.test(String(email).toLowerCase());
    }

    var input = document.querySelector("#phone");
    window.intlTelInput(input, {
        // Seçenekler:
        utilsScript: "https://cdnjs.cloudflare.com/ajax/libs/intl-tel-input/17.0.19/js/utils.js", // Ülke kodları ve formatlama için gerekli
        initialCountry: "tr", // Başlangıç ülkesi (örn: Türkiye)
        separateDialCode: true, // Ülke kodunu ayrı bir alana koyar
    });


    form.addEventListener("submit", function (e) {
        e.preventDefault();

        // ... (mevcut hata kontrolleriniz) ...

        // Input'tan alınan değeri alın, bu artık ülke koduyla birlikte formatlanmış olabilir.
        // Örneğin, "+90 555 555 55 55" gibi.
        // Eğer sadece numarayı backend'e göndermek isterseniz, bunu almanız gerekebilir.
        var phoneInput = document.getElementById("phone");
        var phoneNumber = phoneInput.value; // Bu, kullanıcının girdiği ham değerdir

        // Eğer kütüphanenin biçimlendirilmiş tam numarasını almak isterseniz:
        // var iti = window.intlTelInputGlobals.getInstance(phoneInput);
        // var formattedPhoneNumber = iti.getNumber(); // Formatlanmış E.164 numarası

        // const data = { ... };
        // data.phone = formattedPhoneNumber; // Veya phoneNumber, ne göndermek istediğinize bağlı olarak

        // ... (form gönderme kodlarınız) ...
    });

// Telefon numarası validasyonu için helper fonksiyonunu güncelleyin
    function isValidPhone(phone) {
        var inputElement = document.getElementById("phone");
        var iti = window.intlTelInputGlobals.getInstance(inputElement);
        //intl-tel-input'un kendi isPossibleNumber() ve isValidNumber() metodlarını kullanmak daha iyi olacaktır.
        // Bu metodlar, kullanıcı ülke seçtiğinde ve numara girdiğinde geçerliliği kontrol eder.
        // Eğer sadece basit bir format kontrolü istiyorsanız aşağıdaki gibi kullanabilirsiniz:
        // const re = /^(\d{10}|\d{3}-\d{3}-\d{2}-\d{2})$/; // Bu eski regex'ti
        // return re.test(String(phone));

        //intl-tel-input ile daha iyi validasyon
        return iti.isValidNumber();
    }
});