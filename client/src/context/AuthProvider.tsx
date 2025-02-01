import {
  createContext,
  ReactNode,
  useContext,
  useState,
  useEffect,
} from "react";
import { User } from "../types/types";
import api from "../utils/api";
import { env } from "../utils/env";
import LoadingSpinner from "../components/LoadingSpinner";

interface AuthContextType {
  isAuthenticated: boolean;
  user: User | null;
  login: () => void;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const checkAuthStatus = async () => {
    try {
      const response = await api.get("/auth/user/me");
      setUser(response.data.data);
      setIsAuthenticated(true);
    } catch (error) {
      setIsAuthenticated(false);
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const login = () => {
    window.location.href = `${env.VITE_API_URL}/auth/login`;
  };

  const logout = async () => {
    await api.post("/auth/logout");
    setIsAuthenticated(false);
    setUser(null);
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (ctx === undefined) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
};
