import { create } from "zustand";
import { User } from "../services/types/user";
import { persist } from "zustand/middleware";

interface AuthState {
  accessToken?: string | null;
  user?: User | null;
  setAccessToken: (token: string) => void;
  setUser: (user: User) => void;
  clearAccessToken: () => void;
  clearUser: () => void;
  clearAuthState: () => void;
}

/*
persist: default stores the data in localstorage
*/
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      setAccessToken: (token) => set({ accessToken: token }),
      setUser: (user) => set({ user }),
      clearAccessToken: () => set({ accessToken: null }),
      clearUser: () => set({ user: null }),
      clearAuthState: () => set({ accessToken: null, user: null }),
    }),
    { name: "auth-storage" }
  )
);
