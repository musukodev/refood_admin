package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/idoyudha/duitku-go"
	"github.com/idoyudha/duitku-go/common"
	"github.com/idoyudha/duitku-go/invoice"
	duitkupayment "github.com/idoyudha/duitku-go/payment"
	"github.com/idoyudha/duitku-go/transaction"
	qrcode "github.com/skip2/go-qrcode"

	"refood-backend/config"
)

type DirectV2Request struct {
	MerchantCode    string `json:"merchantCode"`
	PaymentAmount   int    `json:"paymentAmount"`
	PaymentMethod   string `json:"paymentMethod"`
	MerchantOrderId string `json:"merchantOrderId"`
	ProductDetails  string `json:"productDetails"`
	Email           string `json:"email,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	CustomerVaName  string `json:"customerVaName,omitempty"`
	CallbackURL     string `json:"callbackUrl"`
	ReturnURL       string `json:"returnUrl"`
	Signature       string `json:"signature"`
	ExpiryPeriod    int    `json:"expiryPeriod,omitempty"`
}

type DirectV2Response struct {
	MerchantCode  string `json:"merchantCode"`
	Reference     string `json:"reference"`
	PaymentURL    string `json:"paymentUrl"`
	VaNumber      string `json:"vaNumber"`
	Amount        string `json:"amount"`
	QRString      string `json:"qrString"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

type Gateway struct {
	Cfg    *config.Config
	Client *duitku.APIClient
}

func NewGateway(cfg *config.Config) *Gateway {
	client := duitku.NewClient(&common.Config{
		MerchantCode: cfg.MerchantCode,
		APIKey:       cfg.APIKey,
		Environment:  cfg.Environment,
	})

	return &Gateway{Cfg: cfg, Client: client}
}

func (g *Gateway) CreateInvoice(ctx context.Context, order *Order) error {
	req := invoice.CreateInvoiceRequest{
		PaymentAmount:   order.Amount,
		MerchantOrderId: order.ID,
		ProductDetails:  fmt.Sprintf("%s x%d - %s", order.FoodName, order.Qty, order.Store),
		Email:           order.CustomerEmail,
		CustomerVaName:  order.CustomerName,
		PhoneNumber:     order.CustomerPhone,
		CallbackURL:     g.Cfg.CallbackURL(),
		ReturnURL:       g.Cfg.ReturnURL(order.ID),
		ExpiryPeriod:    60,
	}

	log.Printf("[DUITKU] Mengirim CreateInvoice POP untuk Order ID: %s, Amount: %d", order.ID, order.Amount)
	res, httpResp, err := g.Client.InvoiceService.Create(ctx, req)
	if err != nil {
		log.Printf("[ERROR] InvoiceService.Create: %v", err)
		return apiError("gagal membuat invoice Duitku", httpResp, err)
	}
	if httpResp != nil {
		log.Printf("[DUITKU] HTTP Status: %d", httpResp.StatusCode)
	}
	if res.PaymentURL == "" {
		return fmt.Errorf("Duitku tidak mengembalikan payment URL: %s", res.StatusMessage)
	}

	order.Reference = res.Reference
	order.PaymentURL = res.PaymentURL
	order.ResultCode = res.StatusCode
	order.StatusMessage = res.StatusMessage
	return nil
}

func (g *Gateway) CreateDirect(ctx context.Context, order *Order, method common.PaymentMethod) error {
	sigRaw := g.Cfg.MerchantCode + order.ID + strconv.Itoa(order.Amount) + g.Cfg.APIKey
	hash := md5.Sum([]byte(sigRaw))
	signature := hex.EncodeToString(hash[:])

	v2Req := DirectV2Request{
		MerchantCode:    g.Cfg.MerchantCode,
		PaymentAmount:   order.Amount,
		PaymentMethod:   string(method),
		MerchantOrderId: order.ID,
		ProductDetails:  fmt.Sprintf("%s x%d - %s", order.FoodName, order.Qty, order.Store),
		Email:           order.CustomerEmail,
		PhoneNumber:     order.CustomerPhone,
		CustomerVaName:  order.CustomerName,
		CallbackURL:     g.Cfg.CallbackURL(),
		ReturnURL:       g.Cfg.ReturnURL(order.ID),
		Signature:       signature,
		ExpiryPeriod:    60,
	}

	baseURL := common.SandboxV2BaseURL
	if g.Cfg.Environment == common.ProductionEnv {
		baseURL = common.ProductionV2BaseURL
	}

	reqBytes, _ := json.Marshal(v2Req)
	log.Printf("[DUITKU] Mengirim CreateTransaction V2 (%s) untuk Order ID: %s, Amount: %d", method, order.ID, order.Amount)

	var v2Res DirectV2Response
	var apiErrMsg string
	var httpStatus int

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/merchant/v2/inquiry", bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] V2 inquiry POST failed: %v", err)
		apiErrMsg = err.Error()
	} else {
		defer httpResp.Body.Close()
		httpStatus = httpResp.StatusCode
		respBody, _ := io.ReadAll(httpResp.Body)
		log.Printf("[DUITKU] HTTP Status: %d, Response: %s", httpStatus, string(respBody))
		if httpStatus == http.StatusOK {
			json.Unmarshal(respBody, &v2Res)
		} else {
			apiErrMsg = string(respBody)
		}
	}

	vaNumber := v2Res.VaNumber
	qrString := v2Res.QRString
	paymentURL := v2Res.PaymentURL
	reference := v2Res.Reference

	if vaNumber == "" && qrString == "" && paymentURL == "" {
		if g.Cfg.Environment == common.SandboxEnv {
			reference = fmt.Sprintf("MOCK-V2-%d", order.CreatedAt.Unix())
			mStr := string(method)
			if mStr == "BC" || mStr == "M2" || mStr == "I1" {
				prefix := "8808"
				if mStr == "M2" {
					prefix = "8960"
				} else if mStr == "I1" {
					prefix = "9880"
				}
				vaNumber = fmt.Sprintf("%s%08d", prefix, order.CreatedAt.Unix()%100000000)
			} else {
				// QRIS
				qrString = fmt.Sprintf("00020101021226670016ID.CO.DUITKU.WWW011893600000%d0215%s520454115303360540%d5802ID5914REFOOD MARKET6007BATAM6304ABCD", order.CreatedAt.Unix(), order.ID, order.Amount)
			}
		} else {
			return fmt.Errorf("Duitku tidak mengembalikan nomor VA/QRIS untuk metode %s (status %d): %s", method, httpStatus, apiErrMsg)
		}
	}

	order.Reference = reference
	order.PaymentURL = paymentURL
	order.VANumber = vaNumber
	order.QRString = qrString
	order.ResultCode = v2Res.StatusCode
	order.StatusMessage = v2Res.StatusMessage

	if qrString != "" {
		png, err := qrcode.Encode(qrString, qrcode.Medium, 320)
		if err == nil {
			order.QRImage = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		}
	}
	return nil
}

func (g *Gateway) SyncStatus(ctx context.Context, order *Order) error {
	res, _, err := g.Client.TransactionService.GetStatus(ctx, transaction.GetTransactionStatusRequest{
		MerchantOrderId: order.ID,
	})
	if err != nil {
		return err
	}

	order.Reference = res.Reference
	order.ResultCode = res.StatusCode
	order.StatusMessage = res.StatusMessage

	switch res.StatusCode {
	case "00":
		order.Status = StatusSuccess
	case "01":
		order.Status = StatusFailed
	default:
		if order.Status == "" {
			order.Status = StatusPending
		}
	}
	return nil
}

func (g *Gateway) PaymentMethods(ctx context.Context, amount int, datetime string) (duitkupayment.GetPaymentMethodResponse, error) {
	res, _, err := g.Client.PaymentService.GetMethods(ctx, duitkupayment.GetPaymentMethodRequest{
		Amount:   amount,
		Datetime: datetime,
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (g *Gateway) VerifyCallbackSignature(merchantCode, amount, orderID, received string) bool {
	if received == "" {
		return false
	}

	// Signature baru Duitku: HMAC-SHA256(merchantCode + amount + merchantOrderId, apiKey)
	h := hmac.New(sha256.New, []byte(g.Cfg.APIKey))
	h.Write([]byte(merchantCode + amount + orderID))
	hmacSig := hex.EncodeToString(h.Sum(nil))
	if strings.EqualFold(hmacSig, received) {
		return true
	}

	// Fallback signature lama: MD5(merchantCode + amount + merchantOrderId + apiKey)
	m := md5.New()
	m.Write([]byte(merchantCode + amount + orderID + g.Cfg.APIKey))
	md5Sig := hex.EncodeToString(m.Sum(nil))
	return strings.EqualFold(md5Sig, received)
}

func (g *Gateway) SimulationEnabled() bool {
	return g.Cfg.EnableSimulation && g.Cfg.Environment == common.SandboxEnv
}

func apiError(prefix string, httpResp *http.Response, err error) error {
	body := ""
	if httpResp != nil && httpResp.Body != nil {
		if b, readErr := io.ReadAll(httpResp.Body); readErr == nil {
			body = strings.TrimSpace(string(b))
		}
	}

	if body != "" {
		return fmt.Errorf("%s: %v (%s)", prefix, err, body)
	}
	return fmt.Errorf("%s: %v", prefix, err)
}
