import axios from "axios";
import { ENV } from "../../utils/env";
import { useAuthStore } from "../../store/authStore";

/*
axois client for performing actions before and after request
handles access-token refetching
*/
const apiClient = axios.create({
  baseURL: ENV.BASE_URL,
  withCredentials: true,
});

apiClient.interceptors.request.use((config) => {
  const accessToken = useAuthStore.getState().accessToken;
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async function (errors) {
    const originalRequest = errors.config;
    if (errors.response.status === 401 && !originalRequest._retry) {
      try {
        const response = await refreshAuthToken();
        const { access_token } = response.data;
        useAuthStore.getState().accessToken = access_token;
        useAuthStore.getState().isAuthenticated = true;

        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        originalRequest._retry = true;

        return apiClient(originalRequest);
      } catch (error) {
        useAuthStore.getState().logout();
        return Promise.reject(error);
      }
    }
    return Promise.reject(errors);
  }
);

/*
apiClient calls
*/

export const registerUser = async () => {};

export const loginUser = async () => {};

export const refreshAuthToken = async () => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return {
    success: true,
    data: { access_token: "asdf" },
  };
};
