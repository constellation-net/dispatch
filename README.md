# Dispatch

Dispatch is a *very* simple and somewhat primitive ESMTP relay. It will forward incoming email exactly verbatim to the upstream server (for now, anyway).

Rather than relaying each message as it arrives, Dispatch maintains a queue that is periodically flushed to the upstream server in a batch operation. 

You might be asking a very sensible question: why tho? Well, the answer really is "just because I can". I wanted a way to unify my SMTP clients so if I ever had to change the credentials for the upstream server, then I wouldn't have to do that across 50 different apps. There's really no need to have a queue system either, I just felt like adding it for the ✨ ***𝓮𝓯𝓯𝓲𝓬𝓲𝓮𝓷𝓬𝔂*** ✨

## Assumptions & Limitations

As mentioned, this implementation is rather primitive. Currently, it assumes that your upstream server is using [PLAIN SASL authentication](https://datatracker.ietf.org/doc/html/rfc4616) only. Likewise, it also assumes that your clients will be using [PLAIN SASL authentication](https://datatracker.ietf.org/doc/html/rfc4616) only.

## Scalability & High Availability

Dispatch's stateless nature means it can be horizontally scaled infinitely with something like a Kubernetes cluster. Using a DaemonSet and Service is all you really need to achieve High Availability with Dispatch.

## Configuration

Dispatch is configured using environment variables, see [`example.env`](`example.env`) for a list of possible variables.

The `RELAY_PASSWORD` should be a [bcrypt hash](https://bcrypt-generator.com) of the password you want your clients to send to the relay. This is partly for security, but mostly so I have an excuse to use bcrypt.


## Docker Compose

```yaml
services:
  dispatch:
    container_name: dispatch
    image: ghcr.io/constellation-net/dispatch:latest
    ports:
      - "25:25"
    environment:
      DISPATCH_INTERVAL: "60"
      RELAY_HOST: dispatch.starsystem.dev
      RELAY_PORT: "25"
      RELAY_USERNAME: lorem
      RELAY_PASSWORD: bcrypt-hash
      UPSTREAM_HOST: smtp.gmail.com
      UPSTREAM_PORT: "587"
      UPSTREAM_USER: gmail
      UPSTREAM_PASS: password
      UPSTREAM_FROM: noreply@starsystem.dev
      UPSTREAM_REPLYTO: admin@starsystem.dev
    restart: always
```