`Go-smtp-gateway`

## `How to Install`

```
cd deploy
ansible-playbook -i inventory.ini smtp-gateway.yaml
```

## `Test`

```
netcat 192.168.1.100 2525
EHLO localhost
AUTH PLAIN 
AHVzZXJuYW1lAHBhc3N3b3Jk
MAIL FROM:<root@nas.local>
RCPT TO:<info@test.local>
DATA
Subject: Test

Hello from the SMTP test.
.
```
where

```
echo "AHVzZXJuYW1lAHBhc3N3b3Jk" | base64 -d
usernamepassword%     
```

![alt text](image.png)

## `References`

- https://github.com/emersion/go-smtp/tree/master