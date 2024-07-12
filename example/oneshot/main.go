package main //nolint:cyclop // for example

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/gotokatsuya/amazon-pay-sdk-go/amazonpay"
)

// envs.
var (
	publicKeyID    = os.Getenv("AMAZON_PAY_PUBLIC_KEY_ID")
	privateKeyPath = os.Getenv("AMAZON_PAY_PRIVATE_KEY_PATH")
	storeID        = os.Getenv("AMAZON_PAY_STORE_ID")
	merchantID     = os.Getenv("AMAZON_PAY_MERCHANT_ID")
)

const htmlDir = "./example/recurring"

func main() { //nolint:gocognit // for example
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}

	amazonpayCli, err := amazonpay.New(publicKeyID, privateKeyData, "jp", true, http.DefaultClient)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		prescriptionID := uuid.New().String()
		req := &amazonpay.CreateCheckoutSessionRequest{
			WebCheckoutDetails: &amazonpay.WebCheckoutDetails{
				CheckoutReviewReturnURL: fmt.Sprintf("http://localhost:8000/review?prescriptionID=%s", prescriptionID),
			},
			StoreID: storeID,
			Scopes:  []string{"name", "email", "phoneNumber", "shippingAddress"},
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

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		checkoutSessionID := r.URL.Query().Get("amazonCheckoutSessionId")
		prescriptionID := r.URL.Query().Get("prescriptionID")
		resp, httpResp, err := amazonpayCli.GetCheckoutSession(r.Context(), checkoutSessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			data := struct {
				CheckoutSessionID string
				PaymentDescriptor string
				PrescriptionID    string
			}{
				CheckoutSessionID: resp.CheckoutSessionID,
				PaymentDescriptor: resp.PaymentPreferences[0].PaymentDescriptor,
				PrescriptionID:    prescriptionID,
			}
			if err := template.Must(template.ParseFiles(filepath.Join(htmlDir, "review.html"))).Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/approve", func(w http.ResponseWriter, r *http.Request) {
		checkoutSessionID := r.URL.Query().Get("amazonCheckoutSessionId")
		prescriptionID := r.URL.Query().Get("prescriptionID")
		resp, httpResp, err := amazonpayCli.UpdateCheckoutSession(r.Context(), checkoutSessionID, &amazonpay.UpdateCheckoutSessionRequest{
			WebCheckoutDetails: &amazonpay.WebCheckoutDetails{
				CheckoutResultReturnURL: fmt.Sprintf("http://localhost:8000/confirm?prescriptionID=%s", prescriptionID),
			},
			PaymentDetails: &amazonpay.PaymentDetails{
				PaymentIntent:                 "AuthorizeWithCapture",
				CanHandlePendingAuthorization: amazonpay.Bool(false),
				ChargeAmount: &amazonpay.Price{
					Amount:       "1000",
					CurrencyCode: "JPY",
				},
			},
			MerchantMetadata: &amazonpay.MerchantMetadata{
				MerchantReferenceID: prescriptionID,
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer httpResp.Body.Close()
		switch httpResp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			http.Redirect(w, r, resp.WebCheckoutDetails.AmazonPayRedirectURL, http.StatusFound)
		default:
			http.Error(w, resp.ErrorResponse.ReasonCode+" | "+resp.ErrorResponse.Message, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		checkoutSessionID := r.URL.Query().Get("amazonCheckoutSessionId")
		prescriptionID := r.URL.Query().Get("prescriptionID")
		resp, httpResp, err := amazonpayCli.CompleteCheckoutSession(r.Context(), checkoutSessionID, &amazonpay.CompleteCheckoutSessionRequest{
			ChargeAmount: &amazonpay.Price{
				Amount:       "1000",
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
				log.Println("chargeID:", resp.ChargeID)
				log.Println("chargePermissionID:", resp.ChargePermissionID)
				log.Println("MerchantMetadata:", resp.MerchantMetadata)
				log.Println("prescriptionID:", prescriptionID)
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
