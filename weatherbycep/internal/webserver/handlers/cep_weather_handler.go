package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	config "github.com/LucasBelusso1/go-OTELChallange/weatherbycep/configs"
	"github.com/LucasBelusso1/go-OTELChallange/weatherbycep/internal/dto"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const serviceName = "weatherbycep"

func GetTemperatureByZipCode(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	tracer := otel.Tracer("cepvalidation")

	ctx, span := tracer.Start(ctx, serviceName+"-requestCEP")
	cepResponse, err := requestCEP(ctx, cep)
	span.End()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if cepResponse.Localidade == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Can not find zipcode"))
		return
	}

	ctx, span = tracer.Start(ctx, serviceName+"-requestWeather")
	weatherResponse, err := requestWeather(ctx, cepResponse)
	span.End()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if weatherResponse == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	apiResponse := dto.WeatherApiOutput{
		City:  cepResponse.Localidade,
		TempC: weatherResponse.Current.TempC,
		TempF: weatherResponse.Current.TempF,
		TempK: weatherResponse.Current.TempC + 273.15,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)
}

func requestCEP(ctx context.Context, cep string) (dto.ViaCepOutput, error) {
	var viaCepDto dto.ViaCepOutput

	req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return viaCepDto, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return viaCepDto, err
	}

	if res.StatusCode != 200 {
		return viaCepDto, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return viaCepDto, err
	}

	err = json.Unmarshal(body, &viaCepDto)

	if err != nil {
		return viaCepDto, err
	}

	return viaCepDto, nil
}

func requestWeather(ctx context.Context, data dto.ViaCepOutput) (*dto.WeatherOutput, error) {
	var weatherDto *dto.WeatherOutput
	configs := config.GetConfig()

	url := "http://api.weatherapi.com/v1/current.json?key=" + configs.WeatherApiKey + "&q=" + url.QueryEscape(data.Localidade)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if res.StatusCode != 200 {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return weatherDto, err
	}

	err = json.Unmarshal(body, &weatherDto)

	if err != nil {
		return weatherDto, err
	}

	return weatherDto, nil
}
