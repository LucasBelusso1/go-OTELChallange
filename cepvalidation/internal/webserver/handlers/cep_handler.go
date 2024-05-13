package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/LucasBelusso1/go-OTELChallange/cepvalidation/internal/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func ValidateCEPAndDispatch(w http.ResponseWriter, r *http.Request) {
	var cepRequestBody dto.CepRequestBody

	err := json.NewDecoder(r.Body).Decode(&cepRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Invalid zipcode"))
		return
	}

	if len(cepRequestBody.Cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Invalid zipcode"))
		return
	}

	cepRegexCompiled, err := regexp.Compile(`^\d{8}$`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	match := cepRegexCompiled.Match([]byte(cepRequestBody.Cep))

	if !match {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Invalid zipcode"))
		return
	}

	ctx := r.Context()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://weatherbycep:8081/"+cepRequestBody.Cep, nil)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(body)
		return
	}

	var apiOutput dto.ApiOutput
	err = json.NewDecoder(resp.Body).Decode(&apiOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiOutput)
}
