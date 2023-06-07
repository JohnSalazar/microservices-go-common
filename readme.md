# **Microservices in Golang - Common Module**

## This is the common and shared module in the implementation of microservices for the e-commerce application

## Languages

- Golang

---

## Implementations

- Config Parse
- Discovery Service Consul
- GRPC Email Server / Client
- HTTP Server
- HTTP Responses
- Logger
- Metrics Implementation
- Middlewares
  - Authentication
  - Authorization
  - Cors
  - Healthy Check
  - Metrics GRPC
  - Metrics HTTP
- Models
  - ECDSA Public Key
  - Status Payment
  - RSA Public Key
  - Token Claims
- Nats
- Repositories
  - Admin MongoDB Mongo Exporter
- Security
  - Manager Certificates
  - Manager Security ECDSA Keys
  - Manager Security RSA Keys
  - Manager Token
- Services
  - Admin MongoDB Mongo Exporter
  - Certificates
  - Email Service
  - Metrics
  - Security ECDSA Keys
  - Security RSA Keys
- Tasks (background service)
  - Check Certificates
  - Check Service Name (discovery service)
  - Verify Public ECDSA Key (rotation public key)
  - Verify Public RSA Key (rotation public key)
- Tracing
  - Jaeger Provider
  - Open Telemetry Span
- Validators Configuration

---

## Structure Overview

<p align="center">
    <img alt="architecture overview" src="https://github.com/JohnSalazar/microservices-go-common/assets/16736914/3c4b03fe-a7c9-4c94-811e-80310e58c73e" />
</p>

---

## List of Services

- [Authentication](https://github.com/JohnSalazar/microservices-go-authentication)
- [Email](https://github.com/JohnSalazar/microservices-go-email)
- [Customer](https://github.com/JohnSalazar/microservices-go-customer)
- [Product](https://github.com/JohnSalazar/microservices-go-product)
- [Cart](https://github.com/JohnSalazar/microservices-go-cart)
- [Order](https://github.com/JohnSalazar/microservices-go-order)
- [Payment](https://github.com/JohnSalazar/microservices-go-payment)
- [Web](https://github.com/JohnSalazar/microservices-go-web)

---

## About

Common module was developed by [oceano.dev](https://oceano.dev/) <img alt="Brasil" src="https://github.com/JohnSalazar/microservices-go-common/assets/16736914/6bb3c0ad-6fb1-4740-bb25-7ab5bcd83da1" width="20" height="14" /> team under the [MIT license](LICENSE).
