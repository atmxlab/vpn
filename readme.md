## About

This is VPN project for learning

## Init

### Config

```bash
cp .env.example .env 
```

```bash 
cp config/server.example.json config/server.json
```

```bash 
cp config/client.example.json config/client.json
```

### Docker Environment

```bash
make test-env-up
```

### Run server

```bash 
make exec:server
```

```bash 
make run:server
```

### Run client

```bash 
make exec:client
```

```bash 
make run:client
```

## Testing

### Unit & Acceptance

```bash
make test
```

### Manual

#### Check traffic in target host

```bash
make exec:target

tcpdump -p icmp
```

#### Send traffic to target host

```bash
make exec:client

ping 172.16.0.20
```

#### Check

- ping response
- traffic incoming to target
- client and server logs

#### Other test methods

Ping over concrete interface

```bash
ping -I tun0 1.1.1.1
```

Ping over concrete interface with packet len

```bash
ping -M do -I tun0 -s 1300 1.1.1.1
```

HTTP request over concrete interface

```bash
curl --interface tun0 --connect-timeout 3 ip-api.com
```

Check route

```bash
traceroute -i tun0 1.1.1.1
```

Check route using icmp

```bash
traceroute -i tun0 --icmp 1.1.1.1
```

## Action Points for upgrading

- [ ] Develop network settings restore feature
- [ ] Use a library for network settings management instead of use cmd.Exec
- [ ] Write tests for VPN client
- [ ] Develop traffic encrypt feature
