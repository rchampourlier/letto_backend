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

For now, you can have a single workflow, written in JS, that will be run everytime
the `/api/triggers/webhook` endpoint is called.

To edit the workflow, simply edit the `exec/js/main.js` file.

If you have issues when the workflow runs, you can show the traces by getting the
container's logs. First get the last container id using `docker ps -a` and
display its logs with `docker logs <ID>`.

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

We follow the Trello documentation which you will find [here](https://trello.readme.io/v1.0/reference#introduction),
but also provide a summary of instructions to get you on track faster. Please note
it may not be up-to-date at your time of reading, and feel free to open a pull-request
to fix it!

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
