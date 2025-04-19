import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import App from "./App";

// Basic smoke test for rendering login/register form
it("renders login form by default", () => {
  render(<App />);
  expect(screen.getByPlaceholderText(/username/i)).toBeInTheDocument();
  expect(screen.getByPlaceholderText(/password/i)).toBeInTheDocument();
  expect(screen.getByText(/login/i)).toBeInTheDocument();
});

// Switch to register form
it("can switch to register form", () => {
  render(<App />);
  fireEvent.click(screen.getByText(/need an account\? register/i));
  expect(screen.getByText(/register/i)).toBeInTheDocument();
});
