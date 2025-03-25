# ISA-INVESTMENTS-FUNDS

[![Build Status][actions-badge]][actions-url]

## Introduction

### ISA Investment Funds is an API that providers functionality relating to the following:

- Investments a customer can buy within an ISA
- Information relating to a customer's current investment position
- Information relating to a customer's annual tax-free allowance
- Purchase or sale of investments (WIP)

### Scenarios to consider:

- User with £25000 to invest
- Separation of Workplace/Retail customers
- Strictly forbid user from purchasing more than a single product. (Apply in service layer)

### Out of scope:

- Creating of customer account (including KYC)
- Authentication of the requester
- Ability to purchase multiple funds (although the solution is built in such a way to allow this in the future)
- Simplification of the order process (How realistic is it that orders would be executed instantly? Assumption made 
  that orders are likely to be queued?)
- Configure a broker and trigger events for scenarios such as the above. (SMS, Email and push notifications when funds 
  are deposited/withdrawn/traded)

### Enhancements:

- Use OpenAPI for easier integration with clients
- Review choice of HTTP router
- Review choice of Database
- Build out CI
- Secret management
- Review of error handling when publishing events
- More sophisticated logging and error wrapping
- Sorting, filtering and pagination of responses.
- More sophisticated error returns

## How to run: 

Initialise the postgres instance by running the following command. `docker compose up`

We can then run the service via the terminal by inputting the following: `go run cmd/server/main.go`

For direction on how to make requests to the API, please see the postman collection in the `tools` directory. 

Alternatively feel free to: `curl http://localhost:8080/getInvestmentOverview/1`

`SeedDatabase()` populates the db with some initial data so you might want to modify should you consider expanding the functionality.

