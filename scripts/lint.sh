#!/bin/bash

set -e

golint -set_exit_status $(go list ./...)
