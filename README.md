# Prerequisite
```sh
$ cd $GOPATH/src && git clone https://github.com/kmchen/test-fullstack-loyalty
$ cd test-fullstack-loyalty
$ docker-compose up
```


## Loyalty Backend (requires go1.9.1 or above)

```sh
$ cd backend && go build
$ ./backend (server is runningn on localhost:6061)
```

### Documentaion
```sh
$ cd backend
$ godoc -http=:6060
$ Visit http://localhost:6060/pkg/test-fullstack-loyalty/backend/
```

### monitoring

```sh
$ Make sure backend is up and running
$ Visit http://localhost:6061/metrics
$ Please visit http://localhost:6060/pkg/test-fullstack-loyalty/backend/monitoring/ for more informatoin about metrics
```

## Loyalty Frontent

```sh
$ cd frontent && yarn install
$ yarn start
$ Visit http://localhost:4350
```
