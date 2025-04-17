package main

import "testing"

func TestUser(t *testing.T) {
    if users[0].Name != "Alice" {
        t.Error("Expected first user to be Alice")
    }
}
