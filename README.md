# go-OTELChallange

### Run on dev mode:

Firstly you need to get a new **WEATHER_API_KEY**. To do it you need to create a new account on https://www.weatherapi.com/.

Replace the environment variable on `docker-compose.yaml` file, inside `weatherbycep` service configurations.

Then, with docker compose installed, just simply run:

```SHELL
$ docker compose build
```

```SHELL
$ docker compose up -d
```

Then, on any client, call POST `localhost:8080` passing the CEP via JSON. Here's a example:

```JSON
{
	"cep": "04858467"
}
```