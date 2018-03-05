# Plumb - Automated, Interactive ETL for Go

Plumb is a pipeline library for Go. The goal is to build an interactive
pipeline manager that can grow and shrink queues on demand and replicate 
processing stages on demand. Pipelines are a network of **Processes** and
**Channels**. A process is basically a function running in a loop, while 
channels are used to stream arguments into a process and stream outputs 
to the next process. Plumb exposes a REST API to interact with and monitor
these objects in real time.