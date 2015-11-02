oauth2
======

This is a OAuth2 helper library.

This library implements osin storage with [upper.io](https://upper.io) as storage layer. So it supports all storage that upper.io supports (i.e. MySQL, PostgreSQL, SQLite3, MongoDB).

Structs are defined to be as generic as possible. Data layer is generated with [gourd](https://github.com/gourd/gourd) and hence implementing the [gourd's store interface](https://github.com/gourd/kit/store).
