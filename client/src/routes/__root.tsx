import { createRootRoute, Outlet } from "@tanstack/react-router";
import React from "react";
import { Toaster } from "react-hot-toast";

const TanStackRouterDevtools =
  process.env.NODE_ENV === "production"
    ? () => null // Render nothing in production
    : React.lazy(() =>
        // Lazy load in development
        import("@tanstack/router-devtools").then((res) => ({
          default: res.TanStackRouterDevtools,
          // For Embedded Mode
          // default: res.TanStackRouterDevtoolsPanel
        }))
      );

export const Route = createRootRoute({
  component: () => (
    <div>
      <Outlet />
      <Toaster />
      <TanStackRouterDevtools position="bottom-left" />
    </div>
  ),
});
