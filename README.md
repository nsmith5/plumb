<p align="center">
  <img height="150" src="./assets/logo.svg" alt="plumb"/>
</p>

# Plumb - Interactive ETL Pipelines in Go
[![GoDoc][godoc-img]][godoc-url] [![Build Status][travis-img]][travis-url]

Plumb is a pipeline library for Go. The goal is to build an interactive
pipeline manager that can grow and shrink queues on demand and replicate 
processing stages on demand. Pipelines are a network of **Processes** and
**Channels**. A process is basically a function running in a loop, while 
channels are used to stream arguments into a process and stream outputs 
to the next process. Plumb exposes a REST API to interact with and monitor
these objects in real time. 

Example text!

[godoc-url]: https://godoc.org/github.com/nsmith5/plumb
[godoc-img]: https://godoc.org/github.com/nsmith5/plumb?status.svg

[travis-img]: https://travis-ci.org/nsmith5/plumb.svg?branch=master
[travis-url]: https://travis-ci.org/nsmith5/plumb
