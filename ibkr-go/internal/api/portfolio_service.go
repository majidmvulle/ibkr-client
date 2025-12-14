package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/money"
	portfoliov1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1"
	"github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1/portfoliov1connect"
)

const percentageMultiplier = 100

// PortfolioServiceHandler implements the PortfolioService ConnectRPC service.
type PortfolioServiceHandler struct {
	ibkrClient *ibkr.Client
}

// NewPortfolioServiceHandler creates a new PortfolioService handler.
func NewPortfolioServiceHandler(ibkrClient *ibkr.Client) portfoliov1connect.PortfolioServiceHandler {
	return &PortfolioServiceHandler{
		ibkrClient: ibkrClient,
	}
}

// GetPortfolio retrieves portfolio summary for an account.
func (h *PortfolioServiceHandler) GetPortfolio(
	ctx context.Context,
	req *connect.Request[portfoliov1.GetPortfolioRequest],
) (*connect.Response[portfoliov1.GetPortfolioResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Get portfolio positions from IBKR Gateway.
	positions, err := h.ibkrClient.GetPortfolio(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get portfolio: %w", err))
	}

	// Get account summary for total value and cash balance.
	summary, err := h.ibkrClient.GetAccountSummary(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get account summary: %w", err))
	}

	// Map to proto portfolio.
	totalValue, err := money.FromFloat64(summary.NetLiquidation, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert total value: %w", err))
	}

	cashBalance, err := money.FromFloat64(summary.TotalCashValue, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert cash balance: %w", err))
	}

	portfolio := &portfoliov1.Portfolio{
		AccountId:   accountID,
		TotalValue:  totalValue,
		CashBalance: cashBalance,
		Positions:   make([]*portfoliov1.Position, 0, len(positions)),
	}

	// Map positions.
	for i := range positions {
		pos, err := mapIBKRPositionToProto(&positions[i])
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to map position: %w", err))
		}

		portfolio.Positions = append(portfolio.Positions, pos)
	}

	return connect.NewResponse(&portfoliov1.GetPortfolioResponse{
		Portfolio: portfolio,
	}), nil
}

// GetPositions retrieves all positions for an account.
func (h *PortfolioServiceHandler) GetPositions(
	ctx context.Context,
	req *connect.Request[portfoliov1.GetPositionsRequest],
) (*connect.Response[portfoliov1.GetPositionsResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Get portfolio positions from IBKR Gateway.
	positions, err := h.ibkrClient.GetPortfolio(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get positions: %w", err))
	}

	// Map to proto positions.
	protoPositions := make([]*portfoliov1.Position, 0, len(positions))
	for i := range positions {
		pos, err := mapIBKRPositionToProto(&positions[i])
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to map position: %w", err))
		}

		protoPositions = append(protoPositions, pos)
	}

	_ = accountID

	return connect.NewResponse(&portfoliov1.GetPositionsResponse{
		Positions: protoPositions,
	}), nil
}

// GetAccountSummary retrieves account summary information.
func (h *PortfolioServiceHandler) GetAccountSummary(
	ctx context.Context,
	req *connect.Request[portfoliov1.GetAccountSummaryRequest],
) (*connect.Response[portfoliov1.GetAccountSummaryResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Get account summary from IBKR Gateway.
	summary, err := h.ibkrClient.GetAccountSummary(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get account summary: %w", err))
	}

	// Convert all money fields.
	netLiq, err := money.FromFloat64(summary.NetLiquidation, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert net liquidation: %w", err))
	}

	totalCash, err := money.FromFloat64(summary.TotalCashValue, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert total cash: %w", err))
	}

	buyingPower, err := money.FromFloat64(summary.BuyingPower, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert buying power: %w", err))
	}

	equityWithLoan, err := money.FromFloat64(summary.EquityWithLoanValue, summary.Currency)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to convert equity with loan: %w", err))
	}

	// Map to proto account summary.
	protoSummary := &portfoliov1.AccountSummary{
		AccountId:            accountID,
		NetLiquidation:       netLiq,
		TotalCash:            totalCash,
		BuyingPower:          buyingPower,
		EquityWithLoan:       equityWithLoan,
		MaintenanceMarginReq: summary.RegTMargin,
		Currency:             summary.Currency,
	}

	return connect.NewResponse(&portfoliov1.GetAccountSummaryResponse{
		AccountSummary: protoSummary,
	}), nil
}

// Helper functions for mapping IBKR types to proto types.

func mapIBKRPositionToProto(ibkrPos *ibkr.Position) (*portfoliov1.Position, error) {
	marketValue, err := money.FromFloat64(ibkrPos.MktValue, ibkrPos.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to convert market value: %w", err)
	}

	avgCost, err := money.FromFloat64(ibkrPos.AvgCost, ibkrPos.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to convert average cost: %w", err)
	}

	unrealizedPnl, err := money.FromFloat64(ibkrPos.UnrealizedPnl, ibkrPos.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unrealized P&L: %w", err)
	}

	position := &portfoliov1.Position{
		Symbol:               ibkrPos.ContractDesc,
		Quantity:             ibkrPos.Position,
		MarketValue:          marketValue,
		AverageCost:          avgCost,
		UnrealizedPnlPercent: calculatePnLPercent(ibkrPos.UnrealizedPnl, ibkrPos.AvgCost, ibkrPos.Position),
		UnrealizedPnl:        unrealizedPnl,
	}

	return position, nil
}

// calculatePnLPercent calculates the unrealized P&L percentage.
func calculatePnLPercent(unrealizedPnl, avgCost, quantity float64) float64 {
	if avgCost == 0 || quantity == 0 {
		return 0
	}

	totalCost := avgCost * quantity
	if totalCost == 0 {
		return 0
	}

	return (unrealizedPnl / totalCost) * percentageMultiplier
}
