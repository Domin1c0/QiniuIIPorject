package storage

import "testing"

func setupTestStorage(t *testing.T) *Storage {
	store, err := NewStorage("sqlite", ":memory:", nil)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}
	return store
}

func TestAddUser(t *testing.T) {
	store := setupTestStorage(t)

	user := &User{
		Name:     "Cirno",
		Password: "baka⑨baka",
	}
	dbUser, err := store.AddUser(user)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	if dbUser.Id != 1 {
		t.Errorf("expected user ID 1, got %d", dbUser.Id)
	}
}

func TestGet(t *testing.T) {
	store := setupTestStorage(t)

	// Add test user
	user := &User{
		Name:     "Cirno",
		Password: "baka⑨baka",
	}
	dbUser, err := store.AddUser(user)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	// Add test session & message
	session := &Session{
		UserId: dbUser.Id,
	}
	dbSession, err := store.AddSession(*session)
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}
	if dbSession.Id != 1 {
		t.Errorf("expected session ID 1, got %d", dbSession.Id)
	}

	message := &Message{
		Session_id: dbSession.Id,
		Role:       "user",
		Content:    "みんなー！ チルノのさんすう教室⑨周年だよ～！",
	}
	dbMessage, err := store.AddMessage(*message)
	if err != nil {
		t.Fatalf("AddMessage failed: %v", err)
	}
	if dbMessage.Id != 1 {
		t.Errorf("expected message ID 1, got %d", dbMessage.Id)
	}

	// Test Get

	gotSessions, err := store.GetSessionsByUserID(dbUser.Id, 10)
	if err != nil {
		t.Fatalf("GetSessionsByUserID failed: %v", err)
	}
	if len(gotSessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(gotSessions))
	}

	gotMessages, err := store.GetMessagesBySessionID(dbSession.Id)
	if err != nil {
		t.Fatalf("GetMessagesBySessionID failed: %v", err)
	}
	if len(gotMessages) != 1 {
		t.Errorf("expected 1 message, got %d", len(gotMessages))
	}
	if gotMessages[0].Content != message.Content {
		t.Errorf("expected message content %q, got %q", message.Content, gotMessages[0].Content)
	}
}

// TODO more tests
