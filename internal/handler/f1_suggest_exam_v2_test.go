package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type mockBulbasaur struct{}

func (m *mockBulbasaur) CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (bool, string) {
	return false, ""
}
func (m *mockBulbasaur) ChargeCallingLLM(ctx context.Context, code string) bool {
	return false
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
	assert.Contains(t, err.Error(), "invalid input")
}
