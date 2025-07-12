package handler

import (
	"context"
	"darius/internal/errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type mockBulbasaur struct {
	checkCallingLLMResult  string
	checkCallingLLMError   error
	chargeCallingLLMResult bool
}

func (m *mockBulbasaur) CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (string, error) {
	return m.checkCallingLLMResult, m.checkCallingLLMError
}

func (m *mockBulbasaur) ChargeCallingLLM(ctx context.Context, code string) bool {
	return m.chargeCallingLLMResult
}

func Test_checkCanCall_InvalidUidStr(t *testing.T) {
	h := &handler{
		bulbasaur: &mockBulbasaur{},
	}

	// Simulate gRPC metadata with a non-numeric x-user-id
	md := metadata.Pairs("x-user-id", "abc")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := h.checkCanCall(ctx, "any_llm_caller")
	fmt.Println("Error:", err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid input provided")
}

func Test_checkCanCall_EmptyUserId(t *testing.T) {
	h := &handler{
		bulbasaur: &mockBulbasaur{},
	}

	// Simulate gRPC metadata with empty x-user-id
	md := metadata.Pairs("x-user-id", "")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := h.checkCanCall(ctx, "any_llm_caller")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid input provided")
}

func Test_checkCanCall_NoUserIdHeader(t *testing.T) {
	h := &handler{
		bulbasaur: &mockBulbasaur{},
	}

	// Simulate gRPC metadata without x-user-id header
	ctx := context.Background()

	_, err := h.checkCanCall(ctx, "any_llm_caller")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid input provided")
}

func Test_checkCanCall_NotEnoughCredits(t *testing.T) {
	mockBulb := &mockBulbasaur{
		checkCallingLLMError: errors.Error(errors.ErrNotEnoughCredits),
	}
	h := &handler{
		bulbasaur: mockBulb,
	}

	// Simulate gRPC metadata with valid user ID
	md := metadata.Pairs("x-user-id", "123")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := h.checkCanCall(ctx, "f1_suggest_exam")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not enough credits to perform this operation")
}

func Test_checkCanCall_GeneralError(t *testing.T) {
	mockBulb := &mockBulbasaur{
		checkCallingLLMError: errors.Error(errors.ErrGeneral),
	}
	h := &handler{
		bulbasaur: mockBulb,
	}

	// Simulate gRPC metadata with valid user ID
	md := metadata.Pairs("x-user-id", "123")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := h.checkCanCall(ctx, "f1_suggest_exam")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "An unexpected error occurred")
}

func Test_checkCanCall_Success(t *testing.T) {
	mockBulb := &mockBulbasaur{
		checkCallingLLMResult: "transaction_code_123",
		checkCallingLLMError:  nil,
	}
	h := &handler{
		bulbasaur: mockBulb,
	}

	// Simulate gRPC metadata with valid user ID
	md := metadata.Pairs("x-user-id", "123")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	chargeCode, err := h.checkCanCall(ctx, "f1_suggest_exam")
	assert.NoError(t, err)
	assert.Equal(t, "transaction_code_123", chargeCode)
}

func Test_checkCanCall_DifferentLLMCaller(t *testing.T) {
	mockBulb := &mockBulbasaur{
		checkCallingLLMResult: "transaction_code_456",
		checkCallingLLMError:  nil,
	}
	h := &handler{
		bulbasaur: mockBulb,
	}

	// Simulate gRPC metadata with valid user ID
	md := metadata.Pairs("x-user-id", "456")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	chargeCode, err := h.checkCanCall(ctx, "f2_score")
	assert.NoError(t, err)
	assert.Equal(t, "transaction_code_456", chargeCode)
}

func Test_handleErrorWithStatusCode(t *testing.T) {
	h := &handler{}
	ctx := context.Background()

	// Test with network connection error
	err := h.handleErrorWithStatusCode(ctx, fmt.Errorf("network error"), errors.ErrNetworkConnection)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Network connection error")
}
