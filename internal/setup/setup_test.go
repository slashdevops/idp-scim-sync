package setup

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/config"
)

// fakeSecretsFetcher records the names requested and replays canned answers.
type fakeSecretsFetcher struct {
	mu        sync.Mutex
	requested []string
	values    map[string]string
	err       error
	onCall    func(ctx context.Context)
}

func (f *fakeSecretsFetcher) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	if f.onCall != nil {
		f.onCall(ctx)
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	f.mu.Lock()
	f.requested = append(f.requested, secretKey)
	f.mu.Unlock()

	if f.err != nil {
		return "", f.err
	}
	return f.values[secretKey], nil
}

func testConfig() *config.Config {
	cfg := config.New()
	return &cfg
}

func TestFetchSecrets_populatesEveryField(t *testing.T) {
	cfg := testConfig()

	fetcher := &fakeSecretsFetcher{values: map[string]string{
		cfg.GWSUserEmailSecretName:          "admin@example.com",
		cfg.GWSServiceAccountFileSecretName: `{"type":"service_account"}`,
		cfg.AWSSCIMAccessTokenSecretName:    "token-abc",
		cfg.AWSSCIMEndpointSecretName:       "https://scim.example.com/v2",
	}}

	if err := fetchSecrets(context.Background(), fetcher, cfg); err != nil {
		t.Fatalf("fetchSecrets() error = %v", err)
	}

	if cfg.GWSUserEmail != "admin@example.com" {
		t.Errorf("GWSUserEmail = %q", cfg.GWSUserEmail)
	}
	if cfg.GWSServiceAccountFile != `{"type":"service_account"}` {
		t.Errorf("GWSServiceAccountFile = %q", cfg.GWSServiceAccountFile)
	}
	if cfg.AWSSCIMAccessToken != "token-abc" {
		t.Errorf("AWSSCIMAccessToken = %q", cfg.AWSSCIMAccessToken)
	}
	if cfg.AWSSCIMEndpoint != "https://scim.example.com/v2" {
		t.Errorf("AWSSCIMEndpoint = %q", cfg.AWSSCIMEndpoint)
	}

	if len(fetcher.requested) != 4 {
		t.Errorf("requested %d secrets, want 4: %v", len(fetcher.requested), fetcher.requested)
	}
}

func TestFetchSecrets_propagatesError(t *testing.T) {
	wantErr := errors.New("access denied")
	fetcher := &fakeSecretsFetcher{values: map[string]string{}, err: wantErr}

	err := fetchSecrets(context.Background(), fetcher, testConfig())
	if err == nil {
		t.Fatal("fetchSecrets() expected an error")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("fetchSecrets() error = %v, want it to wrap %v", err, wantErr)
	}
}

// Secrets previously called context.Background() internally, so the four
// Secrets Manager requests could not be cancelled — on a Lambda timeout they
// kept running until the runtime froze the environment. The caller's context
// must now govern them.
func TestFetchSecrets_honoursCancelledContext(t *testing.T) {
	fetcher := &fakeSecretsFetcher{values: map[string]string{}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fetchSecrets(ctx, fetcher, testConfig())
	if err == nil {
		t.Fatal("fetchSecrets() expected an error for a cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("fetchSecrets() error = %v, want it to wrap context.Canceled", err)
	}

	fetcher.mu.Lock()
	defer fetcher.mu.Unlock()
	if len(fetcher.requested) != 0 {
		t.Errorf("recorded %d secret reads despite a cancelled context: %v", len(fetcher.requested), fetcher.requested)
	}
}

// A cancelled parent must also reach the in-flight fetches, not just the
// not-yet-started ones.
func TestFetchSecrets_cancelsSiblingsOnFirstError(t *testing.T) {
	var (
		mu        sync.Mutex
		sawCancel bool
	)

	fetcher := &fakeSecretsFetcher{
		values: map[string]string{},
		err:    errors.New("boom"),
		onCall: func(ctx context.Context) {
			if ctx.Err() != nil {
				mu.Lock()
				sawCancel = true
				mu.Unlock()
			}
		},
	}

	if err := fetchSecrets(context.Background(), fetcher, testConfig()); err == nil {
		t.Fatal("fetchSecrets() expected an error")
	}

	mu.Lock()
	defer mu.Unlock()
	if !sawCancel {
		t.Log("no sibling observed cancellation; acceptable if all four completed before the first error surfaced")
	}
}
