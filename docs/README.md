# ℹ️ Attractify Documentation

The following docs will help you to better understand, testdrive and later set up Attractify in production.

## 🔎 About Attractify

The amount of information a user is confronted with during their search for a specific thing is sometimes overwhelming. That's why we think it is essential to personalize the web and in app experience for your users during their journey.

There are different approaches to solve these challenges. You can integrate web tracking, analyze the data in real time and try to predict the user's journey. However, you still need a service that takes over the evaluation and another service that then personalizes the experience for the user on your web site or app.

Yes, this is possible, but we see two problems here:

- The systems need to be extremely well connected.
- In times of GDPR and CCPA, such sensitive data should not reside with a third-party provider.

And these are the reasons why we developed Attractify. We needed a system that would allow us to personalize websites and apps without having to put the data in someone else's hands.

## Contents

- Getting started with the demo
- Setup Attractify in production
- Extend or modify Attractify

## 👍 Getting started with the demo

In order to try Attractify, we have prepared a demo that includes not only the Attractify platform, but also an example web shop. We named it Sportify and it should help you better understand the idea behind Attractify.

To try the demo you need a local Docker host with Docker Compose installed. We provide a docker-compose file that contains all the required services and also provides the sample store.

You just have to type the following command into your terminal. After a short wait, the individual services will start and you can play around with the demo.

```
curl -O https://raw.githubusercontent.com/inovex/attractify/master/docker-compose.demo.yml; docker-compose up
```

Once Attractify is started, you can visit the example shop under [http://localhost:8000](http://localhost:8000)

The Attractify platform is available under [http://localhost:8080](http://localhost:8080). You can use the following credentials to login:

User: `demo@example.com`\
Password: `demo4321`

To give you a small overview of the features and also see what Attractify can do, we have prepared a [video](https://www.youtube.com/watch?v=Z0FM4jD6F0U).

[![](https://img.youtube.com/vi/Z0FM4jD6F0U/sddefault.jpg)](https://www.youtube.com/watch?v=Z0FM4jD6F0U)

## 🖥 Setup Attractify in production

[to be defined]

## 🛠 Extend or modify Attractify

Because Attractify is open source, you can customize and extend it to suit your needs. We use the following technologies.

Backend:
- Go
- PostgreSQL/CockroachDB
- ClickhouseDB
- Kafka

Frontend:
- Vue.js
- Vuex

To make changes to the front- or backend, you should start the databases and Kafka via Docker Compose file `docker-compose.dev.yml` in the repository root. Then you can start the frontend and backend separately.

### Frontend

The frontend can either be launched in a Docker container if you don't have Node.js installed, or you can launch it locally with your existing Node.js installation.

#### In a Docker container...

```
cd frontend
./dev.sh
```

#### ...or locally

```
cd frontend
yarn serve
```

### Backend

The backend can be started with a local Go installation. In the backend there are different services that perform different tasks.

- `server` - This provides the API and delivers the frontend.
- `cron` - Takes care of the regular execution of routine jobs.
- `consumer` - Receives new tracking events and processes them.
- `attractify` - Is a CLI tool to perform certain admin tasks.

The individual commands are stored in the `server/cmd` directory.

To start the server, the dependencies in the form of the databases and the kafka server must first be started. then the server can be started with the following command.

```
go run cmd/server/main.go config.dev.json
```

Each of the individual commands requires a config file containing the connection details for the databases as well as for the Kafka broker.

For the commands `server`, `cron` and `consumer` the config file is simply written directly after the command name. For the Attractify command the config file is specified via a commandline argument `--config config.json`.

There are two config files in the repo:

- `config.dev.json` contains all the settings to connect to the services in the docker compose environment.
- `config.sample.json` is used by the demo application.

`config.sample.json` can also be used later as a basis for a production environment setup.