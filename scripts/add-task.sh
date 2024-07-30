#!/bin/bash

PROCESSING_TIME=5

curl -X "POST" --data '{"Name":"peter","processingTime":5}' http://localhost:8080/tasks