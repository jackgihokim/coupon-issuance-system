# Coupon Issuance System

A simple gRPC-based coupon management service built with Go for a Technical Assignment.

## Overview

This coupon issuance system provides a robust platform for creating campaigns and issuing coupons to users. 
The system is simple but designed with performance in mind, leveraging gRPC for efficient communication.

## Features

- **Campaign Management**
    - Create campaigns with customizable parameters (name, description, start/end dates)
    - Set coupon issuance limits per campaign
    - Retrieve campaign details and status

- **Coupon Issuance**
    - Issue coupons within active campaigns
    - Automatic validation of campaign period and limits
    - Unique coupon ID generation in real time

- **API Architecture**
    - gRPC API with Protocol Buffers (HTTP is available)
    - Clean separation of concerns with handlers and models

## Technology Stack

- **Language**: Go 1.24.1
- **API Framework**: Connect (gRPC & HTTP)
- **Protocol Definition**: Protocol Buffers
