### Structure

```json
{
  "inbounds": [
    {
      "type": "naive",
      "tag": "naive-in",
      
      "listen": "::",
      "listen_port": 443,
      "tcp_fast_open": false,
      "sniff": false,
      "sniff_override_destination": false,
      "domain_strategy": "prefer_ipv6",
      "proxy_protocol": false,
      
      "network": "udp",
      "users": [
        {
          "username": "sekai",
          "password": "password"
        }
      ],
      "tls": {}
    }
  ]
}
```

!!! warning ""

    HTTP3 transport is not included by default, see [Installation](/#installation).

### Naive Fields

#### network

Listen network, one of `tcp` `udp`.

Both if empty.

#### users

==Required==

Naive users.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

### Listen Fields

#### listen

==Required==

Listen address.

#### listen_port

==Required==

Listen port.

#### tcp_fast_open

Enable tcp fast open for listener.

#### sniff

Enable sniffing.

See [Protocol Sniff](/configuration/route/sniff/) for details.

#### sniff_override_destination

Override the connection destination address with the sniffed domain.

If the domain name is invalid (like tor), this will not work.

#### domain_strategy

One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.

If set, the requested domain name will be resolved to IP before routing.

If `sniff_override_destination` is in effect, its value will be taken as a fallback.

#### proxy_protocol

Parse [Proxy Protocol](https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) in the connection header.