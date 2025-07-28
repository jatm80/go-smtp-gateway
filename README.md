`Go-smtp-gateway`

## `How to Install`

```
cd deploy
ansible-playbook -i inventory.ini smtp-gateway.yaml
```

## `Test`

```
netcat gateway.home 2525
EHLO gateway.home
AUTH PLAIN 
AGdhdGV3YXkAZ2F0ZXdheQ==
MAIL FROM:<root@nas.local>
RCPT TO:<info@test.local>
DATA
Subject: Test

Hello from the SMTP test.
.
```
where

```
echo "AGdhdGV3YXkAZ2F0ZXdheQ==" | base64 -d
gatewaygateway
```

![alt text](image.png)

## `References`

- https://github.com/emersion/go-smtp/tree/master