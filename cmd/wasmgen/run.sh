#!/bin/sh 
shopt -s extglob
file="!(*_test).go"
path=""
go run ${file} --code ${path}