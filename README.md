<h1 align="center">Gosane 🧘‍♀️</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A sane and simple Go REST API template. Clone me and edit me to fit your usecase.

## What is Gosane?

Gosane is a cloneable API template to get you up and running quickly. It has made a lot of decisions for you, but easily allows you to swap out the things you don't like.

## Features

| Service | Description |
| --- | --- |
| Auth 🔑 | Social (FB / Google) as well as email based JWT authentication. |
| Database 💽 | Database support using the amazing https://github.com/ent/ent package. |
| Email ✉️ | There's an example AWS SES implementation and an easily extendable interface. |
| Config 🗃 | Simple JSON and environment based configuration. |
| Monitoring 🕵️ | Prometheus handlers for monitoring. |
| Errors 🔦 | Automatic sentry error logging via: https://sentry.io |
| Validation 👮‍♀️ | Validation using an extended version of the https://github.com/go-playground/validator package. |
| Build / Test 💪 | Automatically build and test your code with built in Github pipelines. |
| Server | The underlying server framework is Gin, so you benefit from all the goodness you can find over at: https://github.com/gin-gonic/gin |

## Structure

> Browse the codebase in VS Code here: https://github1s.com/sno6/gosane

Gosane is structured as follows:

### Handlers

Each handler (or endpoint) is grouped and encapsulated in its own folder as can be seen [here](/api/handler). Firstly, you must define the relative path for the handler group in a file such as [this](/api/handler/user/user.go), and then define each endpoint as a separate handler.

### Services & Stores

A handler interacts with your business logic through services, which are aptly defined in [`/service`](/service). These services interact with your database entities (using ent) via [stores](/store). The flow of information should look something like the following:

Handler <-> Services <-> Stores

A store should never be used directly in a handler.

### Internal

Anything that isn't considered business logic should live here. Typically you want to structure these as small modules that you could rip out and run isolated from the rest of the project, if you had to. Examples include, [email](/internal/email), [database](/internal/database), [sentry (error management)](/internal/sentry), etc.

### That's about it, the rest is up to you.

## How to use Gosane

### 1. Clone the project

```sh
git clone git@github.com:sno6/gosane.git
```

### 2. Run the damn thing

```sh
./run.sh
```

> Note that if the above command errors you may need to give the script executive permissions by running: `chmod +x ./run.sh`