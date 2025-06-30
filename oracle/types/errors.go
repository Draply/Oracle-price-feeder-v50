package types

import (
	cosmoserrors "cosmossdk.io/errors"
)

const ModuleName = "price-feeder"

// Price feeder errors
var (
	ErrProviderConnection  = cosmoserrors.Register(ModuleName, 1, "provider connection")
	ErrMissingExchangeRate = cosmoserrors.Register(ModuleName, 2, "missing exchange rate for %s")
	ErrTickerNotFound      = cosmoserrors.Register(ModuleName, 3, "%s failed to get ticker price for %s")
	ErrCandleNotFound      = cosmoserrors.Register(ModuleName, 4, "%s failed to get candle price for %s")

	ErrWebsocketDial  = cosmoserrors.Register(ModuleName, 5, "error connecting to %s websocket: %w")
	ErrWebsocketClose = cosmoserrors.Register(ModuleName, 6, "error closing %s websocket: %w")
	ErrWebsocketSend  = cosmoserrors.Register(ModuleName, 7, "error sending to %s websocket: %w")
	ErrWebsocketRead  = cosmoserrors.Register(ModuleName, 8, "error reading from %s websocket: %w")
)
