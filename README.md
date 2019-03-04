# grpc-crud
A CRUD API based on gRPC and a mongoDB backend 

## Motivation

This is a CRUD based gRPC service that performs the same requirements as the program before. I find this approach more fun than with REST.

## Overview
To get an overview of this project go to *hirepb/hire.proto*, in there you can get all the services performed by the server and the types used by the messages. I'm doing a mix of unary and stream request from the server.

The main source for the server is at *controller/controller.go*, in there are all the handlers declared in the proto file. The client code is at *controller/client_controller.go*, this file controls both the view and the control component for the client.

## Build

Run the makefile included as "make build". Be sure to run "make deps" before building. For changes on the proto file, do "make proto" to re-generate the stub file.

## Installation

Run the makefile included as "make install" to run.

## Run

Start mongoDB and run the server binary with the "-s" to run it as server and without it to run the client.
There are more options that can be displayed with "--help".

## Improvements

- Better decoupling of the DB transaction and the gRPC server handlers to allow for unit tests.
