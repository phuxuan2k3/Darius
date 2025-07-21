package errors

import "fmt"

const (
	ErrInvalidInput       = "Invalid input provided"
	ErrDatabaseConnection = "Failed to connect to the database"
	ErrJSONParsing        = "Error parsing JSON data"
	ErrJSONUnmarshalling  = "Error unmarshalling JSON data"
	ErrLLMGeneration      = "Error generating response from LLM"
	ErrDataHasNotReady    = "Data has not been ready yet, please try again later"
	ErrNetworkConnection  = "Network connection error"
	ErrGeneral            = "An unexpected error occurred"
	ErrNotEnoughCredits   = "Not enough credits to perform this operation"
	ErrChargingFailed     = "Charging failed, please try again later"
)

func Error(err string) error {
	switch err {
	case ErrInvalidInput:
		return fmt.Errorf(ErrInvalidInput)
	case ErrDatabaseConnection:
		return fmt.Errorf(ErrDatabaseConnection)
	case ErrJSONParsing:
		return fmt.Errorf(ErrJSONParsing)
	case ErrLLMGeneration:
		return fmt.Errorf(ErrLLMGeneration)
	case ErrNetworkConnection:
		return fmt.Errorf(ErrNetworkConnection)
	case ErrJSONUnmarshalling:
		return fmt.Errorf(ErrJSONUnmarshalling)
	case ErrNotEnoughCredits:
		return fmt.Errorf(ErrNotEnoughCredits)
	case ErrChargingFailed:
		return fmt.Errorf(ErrChargingFailed)
	case ErrDataHasNotReady:
		return fmt.Errorf(ErrDataHasNotReady)
	default:
		return fmt.Errorf(ErrGeneral)
	}
}

// GetHTTPStatusCode maps error types to appropriate HTTP status codes
func GetHTTPStatusCode(err error) string {
	if err == nil {
		return "200"
	}

	errMsg := err.Error()
	switch errMsg {
	case ErrInvalidInput:
		return "400"
	case ErrNotEnoughCredits:
		return "402"
	case ErrChargingFailed:
		return "402"
	case ErrNetworkConnection:
		return "503"
	case ErrDatabaseConnection:
		return "503"
	case ErrJSONParsing, ErrJSONUnmarshalling:
		return "422"
	case ErrLLMGeneration:
		return "503"
	case ErrDataHasNotReady:
		return "503"
	default:
		return "500"
	}
}
