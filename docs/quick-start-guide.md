# Quick start guide

In this page, you will find step by step cluster initalization, cluster configuration, and some commands to run agents.

### Cluster

Initialize cluster by running the following command;

```
$ cerstore cluster init --name test-cluster
```

Two files will be created: CA `certificate` and `private key`.

```
├── ca.crt
└── ca.key
```


### Server 

Now it's time to create server certificate. Server will use this certificate to enable mTLS in gRPC service. It is important in here to create server certificate with cluster CA. Server and agent will trust each others certificate, since both server and worker certificates are created and signed by same cluster certificate.

```
$ certstore cluster certificate --cacert ../ca.crt --cakey ../ca.key --name certstore-server --type server
```

This will create `server.crt` and `server.key`

```
├── server.crt
└── server.key
```


Now we are ready to start server. We need a configuration file for that, and we can save it to `server.yaml`:

```
listen-port: 10000
tls-ca-cert: "./ca.crt"
tls-server-cert: "./server.crt"
tls-server-cert-key: "./server.key"
certstore:
  services:
    - name: "internal certificate service"
      type: Simple
      args:
        private-key: "$PATH_OF_YOUR_CERT/internal.key"
        certificate: "$PATH_OF_YOUR_CERT/internal.crt"
```

Start server with the following command. Server will start listening on `listen-port: 10000`

```
$ certstore server start --config server.yaml
```


### Agent

Agents also needs a certificate, which should be signed by same cluster CA. Agent will use this certificate to enable mTLS while requesting to server.

```
$ certstore cluster certificate --cacert ../ca.crt --cakey ../ca.key --name certstore-worker --type worker
```

This will create `worker.crt` and `worker.key`

```
├── worker.crt
└── worker.key
```

Worker configuration `worker.yaml`:

```
server-address: "certstore-server:10000"
tls-ca-cert: "./ca.crt"
tls-worker-cert: "./worker.crt"
tls-worker-cert-key: "./worker.key"
pipelines:
  - name: "renew certificate"
    actions:
      - name: issue-certificate
        args:
          issuer: ""internal certificate service""
          common-name: "mywebpage.com"
          sans: "*.mywebpage.com"
      - name: save-certificate
        args:
          certificate-target-path: /tmp/my.crt
          certificate-key-target-path: /tmp/my.key
  - name: should-renew-certificate-pipeline
    actions:
      - name: should-renew-certificate
        args:
          certificate-path: /tmp/my.crt
      - name: run-pipeline
        args:
          pipeline-name: "renew certificate"
jobs:
  - name: "check certificate"
    pipeline: "should-renew-certificate-pipeline"
```

Also add ip address of `certstore-server` to `/etc/hosts`:

```
IP_ADDRESS_OF_SERVER certstore-server
```

Run a pipeline with the following command. This will issue a certificate by requesting server.

```
$ certstore worker runPipeline --config worker.yaml --pipeline renew certificate
```

We can start agent, which agent will schedule and start jobs in this mode.

```
$ certstore worker start --config worker.yaml
```
