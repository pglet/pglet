[![Build status](https://ci.appveyor.com/api/projects/status/1y807vxsh8s6ia2k/branch/master?svg=true)](https://ci.appveyor.com/project/pglet/pglet/branch/master)

# Pglet - Web UI framework for backend developers

Build web apps like a frontend pro in the language you already know. No knowledge of HTML, CSS or JavaScript required.

## What is Pglet

Pglet (*"piglet"*) is a rich user interface (UI) framework for scripts and programs written in any language. [Python](https://pglet.io/docs/tutorials/python), [Bash](https://pglet.io/docs/tutorials/bash), [PowerShell](https://pglet.io/docs/tutorials/powershell) and [Node.js](https://pglet.io/docs/tutorials/node) are already supported and other languages can be easily added via [Pglet protocol](https://pglet.io/docs/reference/protocol).

Pglet renders web UI, so you can easily build web apps with your favorite language. Knowledge of HTML/CSS/JavaScript is not required as you build UI with [controls](https://pglet.io/docs/reference/controls). Pglet controls are built with [Fluent UI React](https://developer.microsoft.com/en-us/fluentui#/controls/web) to ensure your programs look cool and professional.

## Hello world in Bash

Install Pglet helper functions:

```bash
curl -O https://pglet.io/pglet.sh
```

Create `hello.sh` with the following contents:

```bash
. pglet.sh
pglet_page
pglet_add "text value='Hello, world!"
```

Run `sh hello.sh` and in a new browser window you'll get:

<img src="https://pglet.io/img/docs/quickstart-hello-world.png">

Here is a page served by a local instance of Pglet server started in the background on your computer.

Now, add `PGLET_WEB=true` before `pglet_page`:

```bash
. pglet.sh
PGLET_WEB=true pglet_page
pglet_add "text value='Hello, world!"
```

and instantly make your app available on the web!

## Tutorials

* [Python](https://pglet.io/docs/tutorials/python)
* [Bash](https://pglet.io/docs/tutorials/bash)
* [PowerShell](https://pglet.io/docs/tutorials/powershell)
* [Node.js](https://pglet.io/docs/tutorials/node)

## How it works

Pglet UI does not become embedded into your program, but is being served by an out-of-process Pglet server. Application state and control flow logic lives in your persistent-process program while UI changes and events are communicated to Pglet server via IPC-based [protocol](https://pglet.io/docs/reference/protocol). It allows writing web app as a standalone monolith without any knowledge of request/response model, routing, templating or state management. Pglet server can be run locally, self-hosted in your local network or used as a [hosted service](https://pglet.io/docs/pglet-service).

In a classic client-server architecture front-end communicates to a one or more back-end services. Pglet implements an opposite approach where multiple back-end services scattered across internal network behind a firewall and communicate to a centralized Pglet web server, i.e. front-end service, installed in DMZ or [hosted as a service](https://pglet.io/docs/pglet-service). This design gives a number of advantages:

* Secure by design - your internal services and critical data stay behind the firewall and not accessible from the outside world.
* Apps running next to services and data they process - faster/cheaper access and maximum security.
* Zero deployment - run apps on any server in your network or your development machine, no need to deploy apps to a web server.

## Use cases

* Progress visualization for CI/CD workflows, batch jobs and cron tasks 
* Admin interfaces for internal services
* Web dashboards and monitoring
* Status pages
* Executive reporting
* Registration forms and questionnaires
* Intranet self-service kiosks
* Prototype and throw-away apps