# Nordigen API Client

A library provides a client for accessing Nordigen API (https://nordigen.com/).

## Quick Start

Example of the flow for getting access to the account data described
https://nordigen.com/en/account_information_documenation/integration/quickstart_guide/

```go
mySecretID := "c2256760-abc0-49a2-968d-b4cb4cf715d0"
mySecretKey := "88812918b15...93a59239bb7"

api, err := nordigen.New(mySecretID, mySecretKey)
// handle the error

// STEP 1. Choose a bank
institutions := api.Institution().List("DE")

// STEP 2. Choose desired bank from institutions 

// STEP 3. Create an end user agreement
endUserAgreementRequest := &nordigen.CreateAgreementRequest{
	// ...
}
agreement, err := api.EndUserAgreement().Create(endUserAgreementRequest)
// handle the error

// STEP 4. Create a link for a user to confirm (requisition)
requisitionRequest := &CreateRequisitionRequest{
	// ...
}
requisition, err := api.Requisition().Create(requisitionRequest)
// handle the error
// Provide the link from the returned response to a user
// for confirmation

// STEP 5. Get user's accounts (after the confirmation)
requisitionWithAccounts, err := api.Requisition().Get(requisition.ID)
// handle the error

// STEP 6. Access accounts, balances and transactions
account, err := api.Account().Get(accountId)
balances, err := api.Account().Balance(accountId).Get()
details, err := api.Account().Details(accountId).Get()
transactions, err := api.Account().Transaction(accountId).Get(nil, nil)
```

## Concepts

### Client

For accessing Nordigen API 
(https://nordigen.com/en/docs/account-information/integration/parameters-and-responses/)
create an instance of `*Nordigen` type 
```go
n, err := nordigen.New(mySecretID, mySecretKey)
// or 
n := nordigen.MustNew(mySecretID, mySecretKey)
```

Authentication is done by `*Nordigen`  implicitly while 
calling methods that trigger HTTP requests.

### Resources

`*Nordigen` type provides methods for accessing resources
described in the API spec
(https://nordigen.com/en/docs/account-information/integration/parameters-and-responses/)
```go
// Reference to /api/v2/institutions resource
institutionResource := n.Institution()
// Reference to /api/v2/agreements/enduser resource
agreementResource := n.EndUserAgreement()
// Reference to /api/v2/requisitions resource
requisitionResource := n.Requisition()
// Reference to /api/v2/accounts resource
accountResource := n.Account()
```

#### Actions

Resources have CRUD methods that execute HTTP requests.
Note that not all resources has the complete list of CRUD methods.
The availability of methods is dictated by the API spec 
(https://nordigen.com/en/docs/account-information/integration/parameters-and-responses/)

```go
// POST /api/v2/agreements/enduser
agreement, err := n.EndUserAgreement().Create(request)
// GET /api/v2/agreements/enduser
list, err := n.EndUserAgreement().List()
// GET /api/v2/agreements/enduser/{id}
agreement, err := n.EndUserAgreement().Get(id)
// DELETE /api/v2/agreements/enduser/{id}
err := n.EndUserAgreement().Delete(id)
// end user agreement also has Accept() method
// PUT /api/v2/agreements/enduser/{id}/accept
agreement, err := n.EndUserAgreement().Accept(id, data)
```

#### Nested resources

Resources can have nested resources

```go
// Reference to /api/v2/accounts/{id}/transactions resource
transactionResource := n.Account().Transaction(accountId)
// GET /api/v2/accounts/{id}/transactions
transactions, err := transactionResource.Get(from, to)

// GET /api/v2/accounts/{id}/balances
n.Account().Balance(accountId).Get()
```