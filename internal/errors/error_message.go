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
	default:
		return fmt.Errorf(ErrGeneral)
	}
}
