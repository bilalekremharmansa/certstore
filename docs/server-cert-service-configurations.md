# Certificate server configurations

In this page, you will find certificates service configurations for servers.

#### Simple

Simple certificate server creates and signs certificates with given certificate.

```
....
certstore:
  services:
    - name: "certificate service"
      type: Simple
      args:
        private-key: "$private_key_path"
        certificate: "$PATH_OF_YOUR_CERT/internal.crt"
```



##### Certificate authority

Creates CA certificates

```
....
certstore:
  services:
    - name: "certificate authority cert service"
      type: CertificateAuthority
```



#### Let's Encrypt

Issues Let's Encrypt certificates by using [lego](https://github.com/go-acme/lego) library.

> WARNING: It only support windns dns provider at this moment. Create issue for more. 

```
....
certstore:
  services:
    - name: "lets-encrypt-cert-service"
      type: LetsEncrypt
      args:
        private-key: "./acmeuser.key"
        email: "your@mail.com"
        provider: "windns"
```

