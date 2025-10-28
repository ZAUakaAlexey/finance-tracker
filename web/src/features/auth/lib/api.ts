import { AuthResponse, LoginCredentials } from "@/types/auth.types.ts";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
const API_KEY = import.meta.env.VITE_API_KEY || "";

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      throw new Error("Failed to login");
    }

    return await response.json();
  },
};
