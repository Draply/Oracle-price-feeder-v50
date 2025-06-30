package client

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BroadcastTx attempts to generate, sign and broadcast a transaction with the
// given set of messages. It will also simulate gas requirements if necessary.
// It will return an error upon failure.
//
// Note, BroadcastTx is copied from the SDK except it removes a few unnecessary
// things like prompting for confirmation and printing the response. Instead,
// we return the TxResponse.
func BroadcastTx(clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf, err := prepareFactory(clientCtx, txf)

	if err != nil {
		return nil, err
	}

	// _, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println("BroadcastTx2")
	// txf = txf.WithGas(adjusted)

	// TODO: prefer to use CalculateGas but for now we use fixed gas
	const fixedGas uint64 = 200000

	txf = txf.WithGas(fixedGas)
	unsignedTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}
	unsignedTx.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	// unsignedTx.SetFeePayer(clientCtx.GetFeePayerAddress())

	if err = tx.Sign(context.Background(), txf, clientCtx.GetFromName(), unsignedTx, true); err != nil {
		return nil, err
	}
	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		return nil, err
	}

	resp, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return resp, err
	}
	// Log response details
	if resp != nil {
		// Note: We can't use logger here as it's not available in this context
		// The logging will be done in the calling function
	}

	return resp, nil
}

// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()
	// Check if the key exists in the keyring before proceeding
	_, err := clientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		log.Printf("ERROR: Key not found in keyring for address %s: %v", from.String(), err)

		// List available keys for debugging
		keys, listErr := clientCtx.Keyring.List()
		if listErr != nil {
			log.Printf("ERROR: Failed to list keys: %v", listErr)
		} else {
			log.Printf("Available keys in keyring: %+v", keys)
		}

		return txf, fmt.Errorf("key not found in keyring for address %s: %w", from.String(), err)
	}

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		log.Printf("ERROR: AccountRetriever.EnsureExists failed: %v", err)
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			log.Printf("ERROR: GetAccountNumberSequence failed: %v", err)
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
