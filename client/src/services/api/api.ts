import axios from "axios";
import { ENV } from "../../utils/env";
import { useAuthStore } from "../../store/authStore";

/*
axois client for performing actions before and after request
handles access-token refetching
*/
export const api = axios.create({
  baseURL: ENV.BASE_URL,
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  const accessToken = useAuthStore.getState().accessToken;
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async function (errors) {
    const originalRequest = errors.config;
    if (errors.response.status === 401 && !originalRequest._retry) {
      try {
        const response = await api.post("/refresh-token");
        const { accessToken, user } = response.data;
        useAuthStore.getState().setAccessToken(accessToken);
        useAuthStore.getState().setUser(user);

        originalRequest.headers.Authorization = `Bearer ${accessToken}`;
        originalRequest._retry = true;

        return api(originalRequest);
      } catch (error) {
        useAuthStore.getState().clearAuthState();
        return Promise.reject(error);
      }
    }
    return Promise.reject(errors);
  }
);

/*
api calls
*/

export const registerUser = async () => {};

export const loginUser = async () => {};
