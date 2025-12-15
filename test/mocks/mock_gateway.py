#!/usr/bin/env python3
"""
Mock IBKR Client Portal Gateway for testing
Provides mock responses for IBKR API endpoints
"""

from flask import Flask, jsonify, request
import time

app = Flask(__name__)

# Mock data
MOCK_ACCOUNT_ID = "DU123456"
MOCK_ORDERS = []
MOCK_POSITIONS = []

@app.route('/v1/api/tickle', methods=['POST', 'GET'])
def tickle():
    """Health check endpoint"""
    return jsonify({
        "iserver": {"authStatus": {"authenticated": True}},
        "session": str(int(time.time()))
    })

@app.route('/v1/api/iserver/auth/status', methods=['POST'])
def auth_status():
    """Authentication status"""
    return jsonify({
        "authenticated": True,
        "competing": False,
        "connected": True,
        "message": "",
        "MAC": "00:00:00:00:00:00"
    })

@app.route('/v1/api/iserver/reauthenticate', methods=['POST'])
def reauthenticate():
    """Reauthenticate"""
    return jsonify({"message": "triggered"})

@app.route('/v1/api/iserver/accounts', methods=['GET'])
def get_accounts():
    """Get accounts"""
    return jsonify([MOCK_ACCOUNT_ID])

@app.route('/v1/api/iserver/account/<account_id>/orders', methods=['POST'])
def place_order(account_id):
    """Place order"""
    order_data = request.json
    order_id = f"ORDER{len(MOCK_ORDERS) + 1}"

    order = {
        "order_id": order_id,
        "order_status": "Submitted",
        "encrypt_message": "1"
    }
    MOCK_ORDERS.append(order)

    return jsonify({
        "id": order_id,
        "message": ["Order placed successfully"],
        "order_id": order_id,
        "order_status": "Submitted"
    })

@app.route('/v1/api/iserver/account/<account_id>/order/<order_id>', methods=['POST'])
def modify_order(account_id, order_id):
    """Modify order"""
    return jsonify({
        "order_id": order_id,
        "order_status": "Modified"
    })

@app.route('/v1/api/iserver/account/<account_id>/order/<order_id>', methods=['DELETE'])
def cancel_order(account_id, order_id):
    """Cancel order"""
    return jsonify({
        "order_id": order_id,
        "msg": "Request was submitted",
        "conid": 265598,
        "account": account_id
    })

@app.route('/v1/api/iserver/account/orders', methods=['GET'])
def get_live_orders():
    """Get live orders"""
    return jsonify({
        "orders": [
            {
                "acct": MOCK_ACCOUNT_ID,
                "conidex": "265598",
                "conid": 265598,
                "orderId": 1001,
                "cashCcy": "USD",
                "sizeAndFills": "100",
                "orderDesc": "Bought 100",
                "description1": "AAPL",
                "ticker": "AAPL",
                "secType": "STK",
                "listingExchange": "NASDAQ",
                "remainingQuantity": 100.0,
                "filledQuantity": 0.0,
                "totalSize": 100.0,
                "companyName": "APPLE INC",
                "status": "Submitted",
                "order_ref": "QuickTrade",
                "side": "BUY",
                "price": 150.00,
                "bgColor": "#FFFFFF",
                "fgColor": "#000000"
            }
        ],
        "snapshot": True
    })

@app.route('/v1/api/portfolio/<account_id>/positions/0', methods=['GET'])
def get_positions(account_id):
    """Get positions"""
    return jsonify([
        {
            "acctId": account_id,
            "conid": 265598,
            "contractDesc": "AAPL",
            "position": 100.0,
            "mktPrice": 150.00,
            "mktValue": 15000.00,
            "currency": "USD",
            "avgCost": 145.00,
            "avgPrice": 145.00,
            "realizedPnl": 0.00,
            "unrealizedPnl": 500.00,
            "exchs": "NASDAQ",
            "expiry": None,
            "putOrCall": None,
            "multiplier": None,
            "strike": 0.0,
            "exerciseStyle": None,
            "conExchMap": [],
            "assetClass": "STK",
            "undConid": 0
        }
    ])

@app.route('/v1/api/portfolio/<account_id>/summary', methods=['GET'])
def get_account_summary(account_id):
    """Get account summary"""
    return jsonify({
        "accountcode": account_id,
        "accountready": "true",
        "accounttype": "DEMO",
        "cushion": "1",
        "daytradesremaining": "-1",
        "netliquidation": "100000.00",
        "netliquidation-c": "USD",
        "totalcashvalue": "85000.00",
        "totalcashvalue-c": "USD",
        "equity": "100000.00",
        "previousdayequitywithloanvalue": "99500.00"
    })

@app.route('/v1/api/iserver/marketdata/snapshot', methods=['GET'])
def get_market_data():
    """Get market data snapshot"""
    conids = request.args.get('conids', '').split(',')

    snapshots = []
    for conid in conids:
        if conid:
            snapshots.append({
                "conid": int(conid),
                "conidEx": conid,
                "31": "150.00",  # Last price
                "84": "149.50",  # Bid
                "86": "150.50",  # Ask
                "87": "1000000", # Volume
                "70": "152.00",  # High
                "71": "148.00",  # Low
                "7295": "150.00", # Open
                "7296": "150.00", # Close
                "_updated": int(time.time() * 1000)
            })

    return jsonify(snapshots)

@app.route('/v1/api/iserver/marketdata/history', methods=['GET'])
def get_historical_data():
    """Get historical market data"""
    return jsonify({
        "serverId": "1",
        "symbol": "AAPL",
        "text": "APPLE INC",
        "priceFactor": 1,
        "startTime": "20240101-00:00:00",
        "high": "152.00",
        "low": "148.00",
        "timePeriod": "1d",
        "barLength": 300,
        "mdAvailability": "S",
        "mktDataDelay": 0,
        "outsideRth": False,
        "tradingDayDuration": 390,
        "volumeFactor": 1,
        "priceDisplayRule": 1,
        "priceDisplayValue": "2",
        "negativeCapable": False,
        "messageVersion": 2,
        "data": [
            {"t": int(time.time() - 86400) * 1000, "o": 148.0, "c": 150.0, "h": 152.0, "l": 148.0, "v": 1000000},
            {"t": int(time.time()) * 1000, "o": 150.0, "c": 150.5, "h": 151.0, "l": 149.5, "v": 800000}
        ],
        "points": 2,
        "travelTime": 10
    })

@app.route('/v1/api/iserver/secdef/search', methods=['GET'])
def search_contracts():
    """Search for contracts"""
    symbol = request.args.get('symbol', '')

    return jsonify([
        {
            "conid": 265598,
            "companyHeader": "Apple Inc - Common Stock",
            "companyName": "APPLE INC",
            "symbol": symbol.upper(),
            "description": "APPLE INC",
            "restricted": None,
            "fop": None,
            "opt": None,
            "war": None,
            "sections": [
                {
                    "secType": "STK",
                    "months": "",
                    "exchange": "NASDAQ",
                    "legSecType": None
                }
            ]
        }
    ])

if __name__ == '__main__':
    print("Starting Mock IBKR Gateway on port 5555...")
    app.run(host='0.0.0.0', port=5555, debug=False)
