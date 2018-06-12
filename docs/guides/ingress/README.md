---
title: Ingress | Voyager
menu:
  product_voyager_7.1.0:
    identifier: readme-ingress
    name: Readme
    parent: ingress-guides
    weight: -1
product_name: voyager
menu_name: product_voyager_7.1.0
section_menu_id: guides
url: /products/voyager/7.1.0/guides/ingress/
aliases:
  - /products/voyager/7.1.0/guides/ingress/README/
---

# Guides

Guides show you how to use Voyager as a Kubernetes Ingress controller.

- HTTP
  - [Exposing Service via Ingress](/docs/guides/ingress/http/single-service.md)
  - [Virtual Hosting](/docs/guides/ingress/http/virtual-hosting.md)
  - [Supports Loadbalancer Source Range](/docs/guides/ingress/http/source-range.md)
  - [URL and Request Header Re-writing](/docs/guides/ingress/http/rewrite-rules.md)
  - [Enable CORS](/docs/guides/ingress/http/cors.md)
  - [Custom HTTP Port](/docs/guides/ingress/http/custom-http-port.md)
  - [Using External Service as Ingress Backend](/docs/guides/ingress/http/external-svc.md)
  - [HSTS](/docs/guides/ingress/http/hsts.md)
  - [Forward Traffic to StatefulSet Pods](/docs/guides/ingress/http/statefulset-pod.md)
  - [Configure Sticky session to Backends](/docs/guides/ingress/http/sticky-session.md)
  - [Blue Green Deployments using weighted Loadbalancing](/docs/guides/ingress/http/blue-green-deployment.md)
- TLS/SSL
  - [TLS Termination](/docs/guides/ingress/tls/overview.md)
  - [Backend TLS](/docs/guides/ingress/tls/backend-tls.md)
  - [Supports AWS certificate manager](/docs/guides/ingress/tls/aws-cert-manager.md)
- TCP
  - [TCP LoadBalancing](/docs/guides/ingress/tcp/overview.md)
  - [TCP SNI](/docs/guides/ingress/tcp/tcp-sni.md)
- Configuration
  - [Customize generated HAProxy config via BackendRule](/docs/guides/ingress/configuration/backend-rule.md) (can be used for [http rewriting](https://www.haproxy.com/doc/aloha/7.0/haproxy/http_rewriting.html), add [health checks](https://www.haproxy.com/doc/aloha/7.0/haproxy/healthchecks.html), etc.)
  - [Apply Frontend Rules](/docs/guides/ingress/configuration/frontend-rule.md)
  - [Supported Annotations](/docs/guides/ingress/configuration/annotations.md)
  - [Specify NodePort](/docs/guides/ingress/configuration/node-port.md)
  - [Configure global options](/docs/guides/ingress/configuration/default-options.md)
  - [Configure Custom Timeouts for HAProxy](/docs/guides/ingress/configuration/default-timeouts.md)
  - [Using Custom HAProxy Templates](/docs/guides/ingress/configuration/custom-templates.md)
- Security
  - [Configure Basic Auth for HTTP Backends](/docs/guides/ingress/security/basic-auth.md)
  - [Configure External Auth for HTTP Backends](/docs/guides/ingress/security/oauth.md)
  - [TLS Authentication](/docs/guides/ingress/security/tls-auth.md)
- Monitoring
  - [Exposing HAProxy Stats](/docs/guides/ingress/monitoring/stats.md)
- [Scaling Ingress](/docs/guides/ingress/scaling.md)
- [Placement of Ingress Pods](/docs/guides/ingress/pod-placement.md)
