import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import React from "react";
import { Toaster } from "react-hot-toast";

const TanStackRouterDevtools =
  process.env.NODE_ENV === "production"
    ? () => null
    : React.lazy(() =>
        import("@tanstack/router-devtools").then((res) => ({
          default: res.TanStackRouterDevtools,
        }))
      );

export interface AuthContext {
  accessToken: string | null;
  isAuthenticated: boolean;
  refreshAuthToken: () => Promise<void>;
}

export const Route = createRootRouteWithContext<AuthContext>()({
  component: () => {
    return (
      <div>
        <Outlet />
        <Toaster />
        <TanStackRouterDevtools position="bottom-left" />
      </div>
    );
  },
});
