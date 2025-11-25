#!/bin/bash

# Fetch all products
curl localhost:8080/api/product
# Fetch product number 2
curl localhost:8080/api/product/2
# Fetch non-existant product (show response headers)
curl -i localhost:8080/api/product/288
