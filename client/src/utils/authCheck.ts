import { useNavigate } from "@tanstack/react-router";
import { useAuthStore } from "../store/authStore";

export const authCheck = () => {
  const navigate = useNavigate();
  const { accessToken, user } = useAuthStore.getState();
  if (!accessToken || !user) {
    navigate({
      to: "/auth/login",
    });
    return false;
  }
  return true;
};
