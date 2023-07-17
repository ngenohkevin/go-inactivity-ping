# Go Inactivity Ping

Go Inactivity Ping is a command-line tool written in Go that periodically pings a server to keep it active. It helps prevent server spin-down due to inactivity by sending a ping request after a specified interval.

## About the Project

Many hosting services automatically spin down servers that have been inactive for a certain period of time. This can be problematic if you have a server that performs important tasks. Go Inactivity Ping provides a simple solution to keep your server active.

The tool uses the default go net/http library to establish a connection with your server and send a "ping" message every 20 minutes by default. This activity prevents the server from being marked as inactive by the hosting service.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)


## Installation

To use Go Inactivity Ping, you need to have Go installed on your machine. Follow these steps to install and set up the project:

1. Clone the repository:

```shell
git clone https://github.com/ngenohkevin/go-inactivity-ping.git

cd go-inactivity-ping

go build -o go-inactivity-ping

./go-inactivity-ping
```

By default, Go Inactivity Ping pings the server every 20 minutes. You can modify the interval by editing the code in the main function of the main.go file.

The program will continuously ping the server until it is terminated. To stop the program, press Ctrl+C.
