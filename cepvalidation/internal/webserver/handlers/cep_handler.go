package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/LucasBelusso1/go-OTELChallange/cepvalidation/internal/dto"
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

	res, err := http.Get("http://weatherbycep:8081/" + cepRequestBody.Cep)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
		body, err := io.ReadAll(res.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(body)
		return
	}

	var apiOutput dto.ApiOutput
	err = json.NewDecoder(res.Body).Decode(&apiOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiOutput)
}
