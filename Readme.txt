# Uzbek
docker composeni build qilishdan oldin "sudo docker network create GLNetwork" orqali network yaratib olish kerak.
Localhostdagi postgres yoniq bo'lsa undagi :5432 portini o'chirib qo'ying.
NoificaionService/email/sms.go ichiga sms.SetHeader("To", "example@gmail.com") shu yerga o'z emailingizni kiriting bu vaqtinchalik.

.env va docker composeni qo'shamdi, sababi shaxsiy ma'lumotlar mavjud, o'zingiz project kodlarini tekchirib .env qo'shing ularda asosan portlar va token
token secret keyi ishlatilgan, docker composeni ham to'g'irlang.

# English
Before building docker compose, you need to create a network using "sudo docker network create GLNetwork".
If postgres is enabled on localhost, disable port :5432.
In NoificaionService/email/sms.go sms.SetHeader("To", "example@gmail.com") enter your email here this is temporary.

I didn't add .env and docker compose, because there are private data, check the project codes yourself and add .env, they mainly contain ports and tokens
token secret key is used, also fix docker compose.

# Russian
Перед созданием Docker Compose вам необходимо создать сеть, используя «sudo docker network create GLNetwork».
Если postgres включен на локальном хосте, отключите порт: 5432.
В NoificaionService/email/sms.go sms.SetHeader("To", "example@gmail.com") введите здесь свой адрес электронной почты, это временно.

Я не стал добавлять .env и docker compose, потому что там есть приватные данные, сами проверьте коды проектов и добавьте .env, они в основном содержат порты и токены
используется секретный ключ токена, также исправлено создание Docker.