# History

## 2017-10-01

- Connect API's webhook endpoint with JS execution.

## 2017-09-31

- Prepared for JS execution within a Docker container.
- Executing Docker container from Go code using Docker
  Go SDK.

## 2017-09-26

- Implemented one of the _workflow_ endpoints (_create_).
- Started implementation of the _webhook_ trigger endpoint.
- Started implementation of JS execution:
  - Tested [otto](https://github.com/robertkrimen/otto). Does not fit, because it's limited to JS execution and does not include some libraries provided by browsers or NodeJS (such as `XMLHttpRequest`).
  - Will have to test an approach using Docker.

## 2017-09-16

- Started implementation by tweaking the generated scaffold.

## 2017-08-06

- Initialized the project and started with [Goa's Getting Started](https://goa.design/learn/guide/)).
