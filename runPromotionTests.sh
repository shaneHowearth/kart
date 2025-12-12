#!/bin/bash

# Just a simple script to demonstrate that the coupons/promotion code works for
# the given examples.

# Directory holding the promotion code files.
PROMOTION_CODE_DIR="$HOME/Downloads/oolio/"

# Mixed case check
go run cmd/coupons/main.go -p fifTyOff "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Lower case check
go run cmd/coupons/main.go -p happyhrs "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Cache hit check
go run cmd/coupons/main.go -p FIFTYOFF "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Invalid check.
go run cmd/coupons/main.go -p super1001 "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Multiple successful
go run cmd/coupons/main.go -p FIFTYOFF -p happyhrs "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Multiple Unsuccessful
go run cmd/coupons/main.go -p super1001 -p seven "$PROMOTION_CODE_DIR"couponbase{1,2,3}
# Multiple mixed success and fail
go run cmd/coupons/main.go -p FIFTYOFF -p happyhrs -p tomato -p orange "$PROMOTION_CODE_DIR"couponbase{1,2,3}
