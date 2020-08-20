#!/bin/bash
export GOPATH="/home/vagrant/project/project/"
go run src/amqp_server.go &
export FLASK_APP=server.py
export FLASK_ENV=development
flask run --host=0.0.0.0 &