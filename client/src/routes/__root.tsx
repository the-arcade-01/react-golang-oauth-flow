import { createRootRoute, Outlet } from "@tanstack/react-router";
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

export const Route = createRootRoute({
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
