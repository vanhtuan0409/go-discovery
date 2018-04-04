# Service Discovery

### Installation and boostraping

* Install go dependencies and start docker

```
dep init
docker-compose up
```

* Start as much go server as you want

```
PORT=<port> go run server/main.go
```

* Consul and Fabio will automatically discover and load-balancing servers

### Dashboard

* Consul: **htttp://localhost:8500/ui**
* Fabio: **htttp://localhost:9998**
