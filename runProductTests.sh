#!/bin/bash

# Fetch all products
echo "================= All known products ==========================="
curl localhost:8080/api/product
echo "================================================================"
# Fetch product number 2
echo "===================== product ID 2 ============================="
curl localhost:8080/api/product/2
echo "================================================================"
# Fetch non-existant product (show response headers)
echo "===================== non-existant ============================="
curl -i localhost:8080/api/product/288
echo "================================================================"
