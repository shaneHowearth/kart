#!/bin/bash

# Just a simple script to demonstrate that the coupons/promotion code works for
# the given examples.

# Directory holding the promotion code files.
PROMOTION_CODE_DIR="$HOME/Downloads/oolio/"

# Mixed case check
go run cmd/coupons/main.go fifTyOff "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Lower case check
go run cmd/coupons/main.go happyhrs "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Cache hit check
go run cmd/coupons/main.go FIFTYOFF "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Invalid check.
go run cmd/coupons/main.go super1001 "$PROMOTION_CODE_DIR"couponbase{1,2,3}

