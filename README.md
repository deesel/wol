# wol

`wol` is a simple http service that emits [Wake-on-LAN](https://en.wikipedia.org/wiki/Wake-on-LAN) (WoL) messages (a.k.a. magic packets).

You may ask yourself why to complicate such a simple thing as calling [`ether-wake`](https://linux.die.net/man/8/ether-wake) with HTTP API layer? The motivation behind this tool is that it is hard to reliably deliver broadcast packets in some environments like docker. For example if you run [Home Assistant](https://www.home-assistant.io/) in docker container it is hard to use [wake on lan](https://www.home-assistant.io/integrations/wake_on_lan/) integration. In order to send broadcast packets from docker container you need to use host networking, which may not be desirable for each container. One solution could be to isolate WoL functionality into a dedicated single-purpose container connected to host network, which would be operated using HTTP API. And thus `wol` was born...

## How it works

By default it listens on port `8001` and accepts only `POST` method. The magic packet payload is typically broadcast and is encapsulated typically in a UDP datagram (commonly using ports 0, 7 or 9) or directly in ethernet frame using EtherType `0x0842`. These options are configurable and are specified in a body of http request. 

### Ethernet example
```json
{
    "type": "ethernet",
    "mac": "12:34:56:78:90:ab",
    "interface": "eth0"
}
```

### UDP example
```json
{
    "type": "udp",
    "mac": "12:34:56:78:90:ab",
    "ip": "255.255.255.255",
    "port": 9
}
```

Only `mac` value is mandatory, other values are optional. Defaults are:
```json
{
    "type": "udp",
    "ip": "255.255.255.255",
    "port": 9
}
```

## Authorization

Optionally you can configure Bearer token authorization for the service. API keys specified in [configuration](#configuration) are used as authorization tokens.

## Configuration

Most of the configuration options can be supplied either as environment variables or via configuration file. Default configuration file name is `config.yml`, which is searched for in `.` and `/etc/wol` directories. If you like, you can override configuration file name with `-config=<file>` command line argument.

### Configuration file example
```yaml
---
server:
  address: 0.0.0.0
  port: 8001
auth:
  enabled: true
  apiKeys:
  - name: "key 1"
    key: W3i7ST0fFS3YdfooKbVA03up8o17aHy0
  - name: "key 2"
    key: f33CHXl4XGDqNtszdmTiQ6Ya1JXT443c

```

You can also specify configuration parameters as environment variables using following scheme:
```
<section_name>_<parameter>=<value>
```

### Configuration using environment variables
```bash
SERVER_ADDRESS=0.0.0.0
SERVER_PORT=8001
```

By default, service logs at `info` level, for more verbose logging you can use `-log=debug`.

## Usage example

Using configuration file
```bash
docker run --rm -it --net=host -v $PWD/config.yml:/etc/wol/config.yml deesel/wol
```

Using environment variables
```bash
docker run --rm -it --net=host -e SERVER_PORT=8002 deesel/wol
```
