// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"bytes"
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"math/big"
	"time"
)

// StandardFaucet outlines a standard faucet conforming to the standard ruleset.
type StandardFaucet struct {
	Ruleset Ruleset // Faucet ruleset

	AccountsDatabase *accounts.DB // Accounts database
}

/* BEGIN EXPORTED METHODS */

// NewStandardFaucet initializes a new standard faucet.
func NewStandardFaucet(ruleset Ruleset, accountsDB *accounts.DB) *StandardFaucet {
	return &StandardFaucet{
		Ruleset:          ruleset,    // Set ruleset
		AccountsDatabase: accountsDB, // Set DB
	} // Return new faucet
}

// WorkingDB gets a reference to the faucet working database.
func (faucet *StandardFaucet) WorkingDB() *accounts.DB {
	return faucet.AccountsDatabase // Return DB
}

// AccountCanClaim checks if a given account can claim summercash.
func (faucet *StandardFaucet) AccountCanClaim(account *accounts.Account) bool {
	updatedAccount, err := faucet.AccountsDatabase.QueryAccountByUsername(account.Name) // Set to updated account

	if err != nil { // Check for errors
		return false // Cannot claim; does not exist
	}

	account = updatedAccount // Set to updated account reference

	if account.LastFaucetClaimTime.Sub(time.Now()).Hours() < 24 { // Check less than 24 hours since last claim
		return false // Cannot claim
	}

	return true // Can claim
}

// AccountLastClaim gets the time at which an account last claimed from the faucet.
func (faucet *StandardFaucet) AccountLastClaim(account *accounts.Account) time.Time {
	updatedAccount, err := faucet.AccountsDatabase.QueryAccountByUsername(account.Name) // Set to updated account

	if err != nil { // Check for errors
		return time.Date(2017, time.January, 12, 0, 0, 0, 0, time.UTC) // Return despacito music video publish date
	}

	account = updatedAccount // Set to updated account reference

	return account.LastFaucetClaimTime // Return last claim time
}

// AmountCanClaim gets the max amount an account can claim.
// If the account has already claimed in the last 24 hours, zero is returned.
func (faucet *StandardFaucet) AmountCanClaim(account *accounts.Account) *big.Float {
	for _, user := range faucet.Ruleset.BannedUsers() { // Iterate through banned users
		if bytes.Equal(user.Address.Bytes(), account.Address.Bytes()) { // Check is banned user
			return big.NewFloat(0) // Return zero
		}
	}

	if faucet.AccountLastClaim(account).Sub(time.Now()).Hours() < 24 { // Check less than 24 hours since last claim
		return big.NewFloat(0) // Cannot claim
	}

	return faucet.Ruleset.MaximumClaim24hr() // Return max claim 24 hours
}

/* END EXPORTED METHODS */