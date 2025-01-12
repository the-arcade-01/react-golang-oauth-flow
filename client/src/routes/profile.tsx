import { createFileRoute, redirect } from "@tanstack/react-router";
import Profile from "../components/profile";

export const Route = createFileRoute("/profile")({
  component: Profile,
  loader: async ({ context }) => {
    if (!context.isAuthenticated) {
      try {
        await context.refreshAuthToken();
      } catch (error) {
        redirect({
          to: "/auth/login",
          throw: true,
        });
      }
    }
  },
});
