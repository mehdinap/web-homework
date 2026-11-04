<div align="right" dir="rtl">

# تمرین‌های درس مهندسی اینترنت

پاسخ هر تمرین در پوشۀ همان تمرین قرار دارد:

- **تمرین اول:** پوشۀ <code dir="ltr">1/</code> (تمرین <span dir="ltr">Git</span>)
- **تمرین دوم:** پوشۀ <code dir="ltr">2/</code> (تمرین عملی <span dir="ltr">HTTP</span>)

## تمرین دوم: تمرین عملی <span dir="ltr">HTTP</span>

پاسخ دو پرسش اول و توضیح بخش عملی در فایل <code dir="ltr">2/README.md</code> نوشته شده است:

- **بخش ۱:** چگونه می‌توان یک سرویس را فقط روی شبکهٔ داخلی اجرا کرد؟
- **بخش ۲:** فلسفۀ <span dir="ltr">TLS Termination</span>
- **بخش ۳:** پیاده‌سازی پروکسی معکوس با احراز هویت <span dir="ltr">JWT</span>

کد و پیکربندی بخش عملی:

- <code dir="ltr">2/compose.yaml</code> — پروکسی روی پورت ۸۰ و سرویس قدیمی <span dir="ltr">httpbin</span> روی شبکۀ داخلی
- <code dir="ltr">2/proxy/</code> — کد پروکسی به زبان <span dir="ltr">Go</span>
- <code dir="ltr">2/scripts/token.sh</code> — گرفتن توکن برای آزمایش

</div>
