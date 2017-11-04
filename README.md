# Letto

## Status

**Letto is a work-in-progress!!!** It's not usable yet.

## Objective

With Letto, I want to make it easier to build a custom workflow
to trigger automated processes based on the events triggered by
web services.

It's like IFTTT or Zapier, but developer-oriented, where you
would be able to do anything by writing code, but with the
ease-of-use of a dedicated SAAS tool (so no code deployment,
no developing interactions components, no handling authentication...),
you just focus on the workflow you want to create!

## Howto

### Write a workflow

The data (`data.js`, `secrets.js` and workflow scripts) is loaded from the
host directory mapped to `/tmp/data` (see the `docker-compose.yml` file, under
the `web` service).

In the `<dataDir>/workflows` directory, you can provide the JS scripts for your
workflows.

You may use the `exec/js/workflow_example.js` file to see how to build a workflow
script that is runnable.

Some NPM modules are provided in the NodeJS execution environment, you're free
to use them by requiring them. Currently, the following modules are available:

- request
- ovh

### Set credentials

```
cp credentials.js.example credentials.js
```

Then edit `credentials.js` to inject credentials you need from JS. This file is not
committed.

### Integrate with Trello

To integrated with Trello, you will need to:

1. Be able to perform API calls from your workflow, so you will need API credentials.
2. Setup a webhook in Trello to call your Letto endpoint, which is done in Trello
   through the API only.

We follow the Trello documentation which you will find [here](https://trello.readme.io/v1.0/reference#introduction), but also provide a summary of instructions to get you on track faster. Please note it may not be up-to-date at your time of reading, and feel free to open a pull-request to fix it!

**Getting Trello credentials**

You can go on [this page](https://trello.com/app-key) to get your API key
and a personal token.

NB: At this point, Letto is only intended to enable _you_ to build workflows
    for yourself. As such, we do not need to get tokens for other users, so
    using this personal token should be sufficient for now. When we continue
    improving Letto to make it usable by teams, we'll have to rethink this for
    sure!

Once you have your credentials, update the `credentials.js` file using the 
key and token.

**Interact with the API**

To learn how to use the Trello API and using the API in their sandbox, go 
to their [Reference documentation](https://trello.readme.io/reference#membersidboards).

You can use this to setup your webhooks on Letto's endpoint.

NB: Letto's endpoint is: `http[s]://YOUR-DOMAIN:YOUR-PORT/api/triggers/webhook[/YOUR-GROUP]`

## Running

### Prerequisites

**docker-compose**

Install with: 

    sudo curl -L https://github.com/docker/compose/releases/download/1.16.1/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose

### Update the JS data (data, secrets and workflows)

    docker-compose build execjs

### Rebuild and run

    docker-compose build && docker-compose up web

To run the container manually:

    docker run -v /var/run/docker.sock:/var/run/docker.sock -it --rm lettobackend_web:latest

#### Debug the JS execution

```
docker run -it --rm -v "$PWD/../data":/tmp/data -v "$PWD/../traces":/tmp/traces -w /usr/src/app lettobackend_execjs:latest node 
```

### Access private repository in CI/CD

This may help to build the Docker image:

```
# This part enables `go get` to fetch a private repository on gitlab,
# assuming the deploy_key present in the directory is authorized
# for this repo.
RUN mkdir ~/.ssh
RUN mv deploy_key ~/.ssh/id_rsa
RUN mv deploy_key.pub ~/.ssh/id_rsa.pub
RUN echo "Host *\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config
```
