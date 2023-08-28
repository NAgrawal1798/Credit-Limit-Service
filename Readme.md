# Credit Limit Service

This is a Go project that provides a credit limit service for managing account limits and offers.

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [API Endpoints](#api-endpoints)

## Introduction

The Credit Limit Service is designed to manage credit limits and offers for customer accounts. It exposes a RESTful API that allows creating accounts, creating limit offers, updating limit offer statuses, and fetching active limit offers.

## Getting Started

### Prerequisites

- Go (1.14 or later)
- PostgreSQL (with a database named "credit_limit_db")

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/NAgrawal1798/credit-limit-service.git
   cd credit-limit-service
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Initialize the database:

   ```bash
   Create a PostgreSQL database named "credit_limit_db" and run the db/db.go script to initialize the required tables.
   ```

4. Build and run the application:

   ```bash
   go build
   ./credit-limit-service
   ```

## Usage

### API Endpoints

1. Create Account -

   ```bash
   POST /create-account
   ```

   Creates a new customer account.

2. Get Account

   ```bash
   GET /get-account/{account_id}
   ```

   Retrieves account details for the given account ID.

3. Create Limit Offer

   ```bash
   POST /create-limit-offer
   ```

   Creates a new limit offer for an account.

4. List Active Limit Offers

   ```bash
   GET /list-active-limit-offers?account_id={account_id}&active_date={active_date}
   ```

   Fetches active limit offers for the specified account ID and active date.

5. Update Limit Offer Status

   ```bash
   PUT /update-limit-offer-status/{limit_offer_id}/{status}
   ```

   Updates the status of an active or pending limit offer.
