package app

import "testing"

func TestGreetingWithName(t *testing.T) {
	got := Greeting("gopher")
	if got != "hello, gopher" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, gopher")
	}
}

func TestGreetingWithEmptyName(t *testing.T) {
	got := Greeting("   ")
	if got != "hello, world" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, world")
	}
}
