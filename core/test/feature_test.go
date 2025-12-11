package test

import (
	"context"
	"github.com/Alexigbokwe/goNextCore/core"
	"github.com/Alexigbokwe/goNextCore/core/security"
	"github.com/Alexigbokwe/goNextCore/core/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestService struct {
	Value string
}

func NewTestService() *TestService {
	return &TestService{Value: "injected"}
}

func TestDIInvoke(t *testing.T) {
	c := app.NewContainer()
	c.Register(NewTestService())

	// Test function with injection
	fn := func(s *TestService) string {
		return "Hello " + s.Value
	}

	results, err := c.Invoke(fn)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Hello injected", results[0].String())
}

func TestAsyncAwait(t *testing.T) {
	start := time.Now()

	p := utils.Async(func(ctx context.Context) (string, error) {
		select {
		case <-time.After(100 * time.Millisecond):
			return "done", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	})

	res, err := p.Await()
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, "done", res)
	assert.True(t, duration >= 100*time.Millisecond)
}

func TestAsyncCancellation(t *testing.T) {
	p := utils.Async(func(ctx context.Context) (string, error) {
		select {
		case <-time.After(2 * time.Second):
			return "done", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	})

	// Cancel immediately
	p.Cancel()

	_, err := p.Await()
	// Should be error, either "context canceled" or wrapped
	assert.Error(t, err)
}

func TestAsyncTimeout(t *testing.T) {
	p := utils.Async(func(ctx context.Context) (string, error) {
		select {
		case <-time.After(2 * time.Second): // Will take too long
			return "done", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	})

	p.WithTimeout(100 * time.Millisecond)

	_, err := p.Await()
	assert.Error(t, err)
}

func TestPromiseRace(t *testing.T) {
	fast := utils.Async(func(ctx context.Context) (string, error) {
		return "fast", nil
	})
	slow := utils.Async(func(ctx context.Context) (string, error) {
		time.Sleep(1 * time.Second)
		return "slow", nil
	})

	res, err := utils.PromiseRace(fast, slow)
	assert.NoError(t, err)
	assert.Equal(t, "fast", res)
}

func TestJwtService(t *testing.T) {
	jwtService := security.NewJwtService()
	jwtService.SecretKey = "test_secret"

	claims := map[string]interface{}{
		"user_id": "123",
	}

	token, err := jwtService.Sign(claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	verified, err := jwtService.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, "123", verified["user_id"])
}
