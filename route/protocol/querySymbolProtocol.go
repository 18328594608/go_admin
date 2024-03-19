package protocol

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web_app/dao/mysql"
	"net/http"
)

type SymbolInfoHandler struct {
	Data SymbolInfoData
}

type SymbolInfoData struct {
	Symbol string `json:"symbol"`
}

type SymbolInfo struct {
	Symbol       string  `json:"symbol"`
	Group        string  `json:"group"`
	Digit        int     `json:"digit"`
	Point        int     `json:"point"`
	BS           string  `json:"bs"`
	Leverage     int     `json:"leverage"`
	CreateTime   int64   `json:"create_time"`
	Percentage   int     `json:"percentage"`
	Security     string  `json:"security"`
	Currency     string  `db:"currency" json:"currency"`
	ContractSize int     `json:"contract_size"`
	Fee          float64 `json:"fee"`
	SwapLong     float64 `json:"swap_long"`
	SwapShort    float64 `json:"swap_short"`
	BidSpread    int     `json:"bid_spread"`
	AskSpread    int     `json:"ask_spread"`
	MarginCalc   int     `json:"margin_calc"`
	ProfitCalc   int     `json:"profit_calc"`
	SwapCalc     int     `json:"swap_calc"`
	TickSize     float64 `json:"tick_size"`
	TickPrice    int     `json:"tick_price"`
	Monday       string  `json:"monday"`
	Tuesday      string  `json:"tuesday"`
	Wednesday    string  `json:"wednesday"`
	Thursday     string  `json:"thursday"`
	Friday       string  `json:"friday"`
	Comment      string  `json:"comment"`
}

func (si *SymbolInfo) CopyFromSymbol(symbol mysql.Symbol) {
	si.Symbol = symbol.Symbol
	si.Security = symbol.Security
	si.Digit = symbol.Digit
	si.Currency = symbol.Currency
	si.ContractSize = symbol.ContractSize
	si.Percentage = symbol.Percentage
	si.MarginCalc = symbol.MarginCalc
	si.ProfitCalc = symbol.ProfitCalc
	si.SwapCalc = symbol.SwapCalc
	si.TickSize = symbol.TickSize
	si.TickPrice = symbol.TickPrice
	si.Monday = symbol.Monday
	si.Tuesday = symbol.Tuesday
	si.Wednesday = symbol.Wednesday
	si.Thursday = symbol.Thursday
	si.Friday = symbol.Friday
}

func (si *SymbolInfo) CopyFromFee(fee mysql.Fee) {
	si.Symbol = fee.Symbol
	si.Group = fee.Group
	si.Percentage = fee.Percentage
	si.Fee = fee.Fee
	si.SwapLong = fee.SwapLong
	si.SwapShort = fee.SwapShort
	si.BidSpread = fee.BidSpread
	si.AskSpread = fee.AskSpread
}

func (h *SymbolInfoHandler) Handle(context *gin.Context) {
	// 实现"symbolInfo"的具体处理逻辑
	symbols, err := mysql.GetSymbolsFromSymbol(h.Data.Symbol)
	if err != nil {
		zap.L().Error("Failed to get symbols:", zap.Error(err))
		return
	}

	symbolInfo := SymbolInfo{}
	if len(symbols) > 0 {
		firstSymbol := mysql.Symbol{}
		firstSymbol = symbols[0]
		symbolInfo.CopyFromSymbol(firstSymbol)
	}

	fees, err := mysql.GetFeesFromSymbol(h.Data.Symbol)
	if err != nil {
		zap.L().Error("Failed to get fee:", zap.Error(err))
		return
	}
	//获取详细fee数据
	if len(fees) > 0 {
		firstFee := mysql.Fee{}
		firstFee = fees[0]
		symbolInfo.CopyFromFee(firstFee)
	}
	context.JSON(http.StatusOK, gin.H{"cmd": "symbolInfo", "data": symbolInfo})
}
