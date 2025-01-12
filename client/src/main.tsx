import ReactDOM from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import "./index.css";

import { routeTree } from "./routeTree.gen";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useAuthStore } from "./store/authStore";

const router = createRouter({
  routeTree,
  context: {
    accessToken: useAuthStore.getState().accessToken,
    isAuthenticated: useAuthStore.getState().isAuthenticated,
    refreshAuthToken: useAuthStore.getState().refreshAuthToken,
  },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
});

const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);

  const authContext = {
    accessToken: useAuthStore.getState().accessToken,
    isAuthenticated: useAuthStore.getState().isAuthenticated,
    refreshAuthToken: useAuthStore.getState().refreshAuthToken,
  };

  root.render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} context={authContext} />
    </QueryClientProvider>
  );
}
