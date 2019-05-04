// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/transactions"

	"github.com/valyala/fasthttp"

	summercashCommon "github.com/SummerCash/go-summercash/common"
)

/* BEGIN EXPORTED METHODS */

// SetupTransactionsRoutes sets up all the transactions api-related routes.
func (api *JSONHTTPAPI) SetupTransactionsRoutes() error {
	transactionsAPIRoot := "/api/transactions" // Get transactions API root path

	api.Router.POST(fmt.Sprintf("%s/NewTransaction", transactionsAPIRoot), api.NewTransaction) // Set NewTransaction post

	return nil // No error occurred, return nil
}

// NewTransaction handles a NewTransaction request.
func (api *JSONHTTPAPI) NewTransaction(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             // Allow CORS
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	ctx.Response.Header.Set("Content-Type", "application/json")             // Set content type

	var recipient summercashCommon.Address // Init recipient buffer
	var err error                          // Init error buffer

	fmt.Println("test")

	if string(common.GetCtxValue(ctx, "username")) == "faucet" { // Check wants to send from faucet
		logger.Errorf("user with address %s tried to send tx from faucet account", ctx.RemoteAddr().String()) // Log error

		panic(errors.New("cannot send transaction from faucet wallet")) // Panic
	}

	fmt.Println("test")

	if !strings.Contains(string(common.GetCtxValue(ctx, "recipient")), "0x") { // Check is sending to username
		recipientAccount, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "recipient"))) // Query account

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}

		recipient = recipientAccount.Address // Set address
	} else {
		recipient, err = summercashCommon.StringToAddress(string(common.GetCtxValue(ctx, "recipient"))) // Parse recipient

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}
	}

	fmt.Println(string(common.GetCtxValue(ctx, "amount")))

	amount, err := strconv.ParseFloat(string(common.GetCtxValue(ctx, "amount")), 64) // Parse amount

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Println("test")

	transaction, err := transactions.NewTransaction(api.AccountsDatabase, string(common.GetCtxValue(ctx, "username")), string(common.GetCtxValue(ctx, "password")), &recipient, amount, common.GetCtxValue(ctx, "payload")) // Initialize transaction

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Println("test")

	fmt.Fprintf(ctx, transaction.String()) // Write tx string value
}

/* END EXPORTED METHODS */
