package main //nolint:cyclop // for example

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/sasada-t/amazon-pay-sdk-go/amazonpay"
)

// envs.
var (
	publicKeyID = os.Getenv("AMAZON_PAY_PUBLIC_KEY_ID")
	privateKey  = os.Getenv("AMAZON_PAY_PRIVATE_KEY")
	storeID     = os.Getenv("AMAZON_PAY_STORE_ID")
	merchantID  = os.Getenv("AMAZON_PAY_MERCHANT_ID")
)

// local datastore.
var (
	chargePermissionID string
)

const htmlDir = "./example/recurring"

var refID = uuid.New().String()

func main() { //nolint:gocognit,cyclop // for example
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	amazonpayCli, err := amazonpay.New(publicKeyID, decodedPrivateKey, "jp", true, http.DefaultClient)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req := &amazonpay.CreateCheckoutSessionRequest{
			WebCheckoutDetails: &amazonpay.WebCheckoutDetails{
				CheckoutReviewReturnURL: "http://localhost:8000/approve",
			},
			StoreID:              storeID,
			ChargePermissionType: "Recurring",
			RecurringMetadata: &amazonpay.RecurringMetadata{
				Frequency: &amazonpay.Frequency{
					Unit:  "Variable",
					Value: "0",
				},
				Amount: &amazonpay.Price{
					Amount:       "1",
					CurrencyCode: "JPY",
				},
			},
			PaymentDetails: &amazonpay.PaymentDetails{
				PaymentIntent:                 "Confirm",
				CanHandlePendingAuthorization: amazonpay.Bool(false),
				ChargeAmount: &amazonpay.Price{
					Amount:       "1",
					CurrencyCode: "JPY",
				},
			},
			MerchantMetadata: &amazonpay.MerchantMetadata{
				NoteToBuyer: "Testing plan",
			},
			ProviderMetadata: &amazonpay.ProviderMetadata{
				ProviderReferenceID: refID,
			},
		}
		payload, err := req.ToPayload()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signature, err := amazonpayCli.GenerateButtonSignature(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			AmazonPayPayload     string
			AmazonPaySignature   string
			AmazonPayPublicKeyID string
			AmazonPayMerchantID  string
		}{
			AmazonPayPayload:     payload,
			AmazonPaySignature:   signature,
			AmazonPayPublicKeyID: publicKeyID,
			AmazonPayMerchantID:  merchantID,
		}
		if err := template.Must(template.ParseFiles(filepath.Join(htmlDir, "index.html"))).Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/approve", func(w http.ResponseWriter, r *http.Request) {
		checkoutSessionID := r.URL.Query().Get("amazonCheckoutSessionId")
		resp, httpResp, err := amazonpayCli.UpdateCheckoutSession(r.Context(), checkoutSessionID, &amazonpay.UpdateCheckoutSessionRequest{
			WebCheckoutDetails: &amazonpay.WebCheckoutDetails{
				CheckoutResultReturnURL: "http://localhost:8000/completed",
			},
			PaymentDetails: &amazonpay.PaymentDetails{
				PaymentIntent:                 "Confirm",
				CanHandlePendingAuthorization: amazonpay.Bool(false),
				ChargeAmount: &amazonpay.Price{
					Amount:       "1",
					CurrencyCode: "JPY",
				},
			},
			MerchantMetadata: &amazonpay.MerchantMetadata{
				NoteToBuyer: "Testing plan",
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("approve", resp.WebCheckoutDetails.AmazonPayRedirectURL)
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			http.Redirect(w, r, resp.WebCheckoutDetails.AmazonPayRedirectURL, http.StatusFound)
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/completed", func(w http.ResponseWriter, r *http.Request) {
		checkoutSessionID := r.URL.Query().Get("amazonCheckoutSessionId")
		resp, httpResp, err := amazonpayCli.CompleteCheckoutSession(r.Context(), checkoutSessionID, &amazonpay.CompleteCheckoutSessionRequest{
			ChargeAmount: &amazonpay.Price{
				Amount:       "1",
				CurrencyCode: "JPY",
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			log.Println("confirm: " + resp.StatusDetails.State)
			switch resp.StatusDetails.State {
			case "Open":
			case "Completed":
				// TODO should save to database
				log.Println("ChargeID:", resp.ChargeID)
				log.Println("ChargePermissionID:", resp.ChargePermissionID)
				log.Println("ChargePermissionType:", resp.ChargePermissionType)
				// jsonとしてprintする
				b, _ := json.Marshal(resp)
				log.Println(string(b))
				chargePermissionID = resp.ChargePermissionID
			case "Canceled":
			}
			data := struct{}{}
			if err := template.Must(template.ParseFiles(filepath.Join(htmlDir, "confirm.html"))).Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/recurring", func(w http.ResponseWriter, r *http.Request) {
		refID := uuid.New().String()
		cpResp, httpResp, err := amazonpayCli.GetChargePermission(r.Context(), chargePermissionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			log.Println("recurring: " + cpResp.StatusDetails.State)
			switch cpResp.StatusDetails.State {
			case "Chargeable":
				cResp, httpResp, err := amazonpayCli.CreateCharge(r.Context(), &amazonpay.CreateChargeRequest{
					ChargePermissionID: chargePermissionID,
					ChargeAmount: &amazonpay.Price{
						Amount:       "10000",
						CurrencyCode: "JPY",
					},
					CaptureNow:                    amazonpay.Bool(true),
					CanHandlePendingAuthorization: amazonpay.Bool(false),
					MerchantMetadata: &amazonpay.MerchantMetadata{
						MerchantReferenceID: refID,
					},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer httpResp.Body.Close()
				log.Println(httpResp.StatusCode)
				switch httpResp.StatusCode {
				case http.StatusOK, http.StatusCreated:
					log.Println("Success /recurring")
					w.WriteHeader(http.StatusOK)
					data := struct {
						ChargeID           string
						ChargePermissionID string
						ChargeState        string
					}{
						ChargeID:           cResp.ChargeID,
						ChargePermissionID: chargePermissionID,
						ChargeState:        cResp.StatusDetails.State,
					}
					if err := template.Must(template.ParseFiles(filepath.Join(htmlDir, "recurring.html"))).Execute(w, data); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				default:
					http.Error(w, cResp.ErrorResponse.ReasonCode+" | "+cResp.ErrorResponse.Message, http.StatusInternalServerError)
				}
			case "NonChargeable":
			case "Closed":
			}
		default:
			http.Error(w, cpResp.ErrorResponse.ReasonCode+" | "+cpResp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/charge", func(w http.ResponseWriter, r *http.Request) {
		chargeID := r.URL.Query().Get("chargeID")
		resp, httpResp, err := amazonpayCli.GetCharge(r.Context(), chargeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			log.Println("Success /charge/:chargeID")
			w.WriteHeader(http.StatusOK)
			data := struct {
				ChargeID           string
				ChargePermissionID string
				ChargeState        string
				RefID              string
			}{
				ChargeID:           resp.ChargeID,
				ChargePermissionID: resp.ChargePermissionID,
				ChargeState:        resp.StatusDetails.State,
				RefID:              resp.MerchantMetadata.MerchantReferenceID,
			}
			if err := template.Must(template.ParseFiles(filepath.Join(htmlDir, "charge.html"))).Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/recurring/close", func(w http.ResponseWriter, r *http.Request) {
		resp, httpResp, err := amazonpayCli.CloseChargePermission(r.Context(), chargePermissionID, &amazonpay.CloseChargePermissionRequest{
			ClosureReason:        "closing reason",
			CancelPendingCharges: amazonpay.Bool(false),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			log.Println("Success /recurring/close")
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	fmt.Println("http://localhost:8000")
	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      nil,
	}
	log.Fatalln(server.ListenAndServe())
}
