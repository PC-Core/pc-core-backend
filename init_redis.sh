#!/bin/bash

redis-cli set userID 0
redis-cli expire userID 31556926