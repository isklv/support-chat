package support_test

import (
	"context"
	"path/filepath"
	"testing"

	"support-chat/internal/config"
	"support-chat/internal/db"
	"support-chat/internal/mailer"
	"support-chat/internal/support"
)

func setupTestService(t *testing.T) *support.Service {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := db.Init("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	cfg := &config.Config{
		AppURL: "http://localhost:8080",
	}
	mailClient := mailer.New(cfg)
	return support.NewService(database, cfg, mailClient, nil)
}

func TestSupportService_PostAndGet(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Post first message as guest
	input := support.PostInput{
		Text:  "Hello, I need help!",
		Email: "user@example.com",
		Name:  "Alex",
		Page:  "/pricing",
		Lang:  "ru",
	}

	view, err := svc.Post(ctx, support.Visitor{}, input)
	if err != nil {
		t.Fatalf("post message failed: %v", err)
	}

	if view.Thread == nil {
		t.Fatal("expected thread to be created")
	}
	if view.Thread.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got '%s'", view.Thread.Email)
	}
	if view.Token == "" {
		t.Fatal("expected guest token to be returned")
	}
	if len(view.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(view.Messages))
	}
	if view.Messages[0].Text != "Hello, I need help!" {
		t.Errorf("unexpected message text: %s", view.Messages[0].Text)
	}

	// Verify current view using the returned token
	curView, err := svc.Current(ctx, support.Visitor{Token: view.Token})
	if err != nil {
		t.Fatalf("current failed: %v", err)
	}
	if curView.Thread == nil || curView.Thread.ID != view.Thread.ID {
		t.Errorf("expected to find same thread with token")
	}

	// Operator adds reply
	_, err = svc.AddOperatorMessage(ctx, curView.Thread, "Support Lead", "We are here to help!", nil, "", nil)
	if err != nil {
		t.Fatalf("add operator message failed: %v", err)
	}

	// Check unread count for guest
	unread, err := svc.Unread(ctx, support.Visitor{Token: view.Token})
	if err != nil {
		t.Fatalf("unread failed: %v", err)
	}
	if unread != 1 {
		t.Errorf("expected unread count 1, got %d", unread)
	}

	// Mark seen
	err = svc.MarkSeen(ctx, support.Visitor{Token: view.Token})
	if err != nil {
		t.Fatalf("mark seen failed: %v", err)
	}

	unreadAfter, err := svc.Unread(ctx, support.Visitor{Token: view.Token})
	if err != nil {
		t.Fatalf("unread after mark seen failed: %v", err)
	}
	if unreadAfter != 0 {
		t.Errorf("expected unread count 0, got %d", unreadAfter)
	}
}

func TestSupportService_CloseAndReopen(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	view, err := svc.Post(ctx, support.Visitor{}, support.PostInput{
		Text:  "Question about export",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}

	// Close thread
	closedView, err := svc.SetStatus(ctx, view.Thread.ID, false, "Admin")
	if err != nil {
		t.Fatalf("close status failed: %v", err)
	}
	if closedView.Thread.Status != db.SupportClosed {
		t.Errorf("expected status 'closed', got '%s'", closedView.Thread.Status)
	}

	// Reopen thread
	openView, err := svc.SetStatus(ctx, view.Thread.ID, true, "Admin")
	if err != nil {
		t.Fatalf("reopen status failed: %v", err)
	}
	if openView.Thread.Status != db.SupportOpen {
		t.Errorf("expected status 'open', got '%s'", openView.Thread.Status)
	}
}
