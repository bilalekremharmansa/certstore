# certstore

certstore is an open-source automated certificate management tool. 

Server is a centralized authority, and issues certificates through certificate services.

Agent is an application that requests it's server to issue certificates. It could be run on various machines, and runs customizable pipelines.



### Features

- Issue certificates
- Renew certificates
- Schedule jobs to issue/renew certificates
- Server and Agents communicates over mTLS
- Customizable pipeline and actions

##### Implemented certificate services

- Certificate Authority service
- Simple service - creates certificates with given CA
- Let's Encrypt service

For more, [see](./docs/server-cert-service-configurations.md).

### Usage

Server and Agents can be configured to run seperately. Please, check [quick start guide](./docs/quick-start-guide.md).



### Build

certstore can be built with make command:

```
make build
```

Two executable binary will be built for windows and linux.

Pre-build releases could be found in releases pages.

