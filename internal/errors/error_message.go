package errors

import "fmt"

const (
	ErrInvalidInput       = "Invalid input provided"
	ErrDatabaseConnection = "Failed to connect to the database"
	ErrJSONParsing        = "Error parsing JSON data"
	ErrJSONUnmarshalling  = "Error unmarshalling JSON data"
	ErrLLMGeneration      = "Error generating response from LLM"
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
	default:
		return fmt.Errorf(ErrGeneral)
	}
}
