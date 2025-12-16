package api

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	marketdatav1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1"
	"github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1/marketdatav1connect"
)

const (
	streamInterval = 5 * time.Second
	maxStreamCount = 100
)

// MarketDataServiceHandler implements the MarketDataService ConnectRPC service.
type MarketDataServiceHandler struct {
	ibkrClient ibkr.MarketDataClient
}

// NewMarketDataServiceHandler creates a new MarketDataService handler.
func NewMarketDataServiceHandler(ibkrClient ibkr.MarketDataClient) marketdatav1connect.MarketDataServiceHandler {
	return &MarketDataServiceHandler{
		ibkrClient: ibkrClient,
	}
}

// GetQuote retrieves a market data quote for a symbol.
func (h *MarketDataServiceHandler) GetQuote(
	ctx context.Context,
	req *connect.Request[marketdatav1.GetQuoteRequest],
) (*connect.Response[marketdatav1.GetQuoteResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Search for contract by symbol.
	contracts, err := h.ibkrClient.SearchContracts(ctx, req.Msg.Symbol)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to search contracts: %w", err))
	}

	if len(contracts) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("symbol not found: %s", req.Msg.Symbol))
	}

	// Use the first contract found.
	conID := contracts[0].ConID

	// Get market data snapshot.
	snapshots, err := h.ibkrClient.GetMarketData(ctx, []int{conID}, nil)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get market data: %w", err))
	}

	if len(snapshots) == 0 {
		return nil, connect.NewError(
			connect.CodeNotFound,
			fmt.Errorf("no market data available for symbol: %s", req.Msg.Symbol),
		)
	}

	// Map to proto quote.
	quote := mapSnapshotToQuote(&snapshots[0], req.Msg.Symbol)

	_ = accountID

	return connect.NewResponse(&marketdatav1.GetQuoteResponse{
		Quote: quote,
	}), nil
}

// GetHistoricalData retrieves historical market data for a symbol.
func (h *MarketDataServiceHandler) GetHistoricalData(
	ctx context.Context,
	req *connect.Request[marketdatav1.GetHistoricalDataRequest],
) (*connect.Response[marketdatav1.GetHistoricalDataResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Search for contract by symbol.
	contracts, err := h.ibkrClient.SearchContracts(ctx, req.Msg.Symbol)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to search contracts: %w", err))
	}

	if len(contracts) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("symbol not found: %s", req.Msg.Symbol))
	}

	// Use the first contract found.
	conID := contracts[0].ConID

	// Get historical data.
	histData, err := h.ibkrClient.GetHistoricalData(ctx, conID, req.Msg.Period, req.Msg.BarSize)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get historical data: %w", err))
	}

	// Map to proto bars.
	bars := make([]*marketdatav1.Bar, 0, len(histData.Data))
	for i := range histData.Data {
		bar := mapHistoricalBarToProto(&histData.Data[i])
		bars = append(bars, bar)
	}

	_ = accountID

	return connect.NewResponse(&marketdatav1.GetHistoricalDataResponse{
		Bars: bars,
	}), nil
}

// StreamQuotes streams real-time quotes for a symbol.
func (h *MarketDataServiceHandler) StreamQuotes(
	ctx context.Context,
	req *connect.Request[marketdatav1.StreamQuotesRequest],
	stream *connect.ServerStream[marketdatav1.StreamQuotesResponse],
) error {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	_ = accountID

	// Search for contract by symbol.
	contracts, err := h.ibkrClient.SearchContracts(ctx, req.Msg.Symbol)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to search contracts: %w", err))
	}

	if len(contracts) == 0 {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("symbol not found: %s", req.Msg.Symbol))
	}

	// Use the first contract found.
	conID := contracts[0].ConID

	return h.streamQuotesLoop(ctx, conID, req.Msg.Symbol, stream)
}

// streamQuotesLoop handles the streaming loop for quotes.
func (h *MarketDataServiceHandler) streamQuotesLoop(
	ctx context.Context,
	conID int,
	symbol string,
	stream *connect.ServerStream[marketdatav1.StreamQuotesResponse],
) error {
	ticker := time.NewTicker(streamInterval)
	defer ticker.Stop()

	count := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := h.sendQuoteSnapshot(ctx, conID, symbol, stream); err != nil {
				return err
			}

			count++
			if count >= maxStreamCount {
				return nil
			}
		}
	}
}

// sendQuoteSnapshot fetches and sends a single quote snapshot.
func (h *MarketDataServiceHandler) sendQuoteSnapshot(
	ctx context.Context,
	conID int,
	symbol string,
	stream *connect.ServerStream[marketdatav1.StreamQuotesResponse],
) error {
	snapshots, err := h.ibkrClient.GetMarketData(ctx, []int{conID}, nil)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get market data: %w", err))
	}

	if len(snapshots) == 0 {
		return nil
	}

	quote := mapSnapshotToQuote(&snapshots[0], symbol)

	return stream.Send(&marketdatav1.StreamQuotesResponse{
		Quote: quote,
	})
}

// Helper functions for mapping IBKR types to proto types.

func mapSnapshotToQuote(snapshot *ibkr.MarketDataSnapshot, symbol string) *marketdatav1.Quote {
	quote := &marketdatav1.Quote{
		Symbol: symbol,
		Last:   snapshot.LastPrice,
		Bid:    snapshot.Bid,
		Ask:    snapshot.Ask,
		Volume: snapshot.Volume,
		High:   snapshot.High,
		Low:    snapshot.Low,
	}

	return quote
}

func mapHistoricalBarToProto(bar *ibkr.HistoricalBar) *marketdatav1.Bar {
	protoBar := &marketdatav1.Bar{
		Timestamp: fmt.Sprintf("%d", bar.Time),
		Open:      bar.Open,
		High:      bar.High,
		Low:       bar.Low,
		Close:     bar.Close,
		Volume:    bar.Volume,
	}

	return protoBar
}
