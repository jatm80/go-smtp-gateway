`Go-smtp-gateway`

## `How to Install`

## `References`

## `Test`

```
netcat 192.168.1.100 2525
EHLO localhost
AUTH PLAIN
AHVzZXJuYW1lAHBhc3N3b3Jk
MAIL FROM:<root@nas.local>
RCPT TO:<info@test.local>
DATA
.
```