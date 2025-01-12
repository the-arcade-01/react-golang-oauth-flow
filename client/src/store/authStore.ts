import { create } from "zustand";
import { refreshAuthToken } from "../services/api/api";

interface AuthState {
  accessToken: string | null;
  isAuthenticated: boolean;
  refreshAuthToken: () => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  isAuthenticated: false,
  refreshAuthToken: async () => {
    try {
      const response = await refreshAuthToken();
      set({
        accessToken: response.data.access_token,
        isAuthenticated: true,
      });
    } catch (error) {
      set({ accessToken: null, isAuthenticated: false });
      throw error;
    }
  },
  logout: () => {
    set({ accessToken: null, isAuthenticated: false });
  },
}));
