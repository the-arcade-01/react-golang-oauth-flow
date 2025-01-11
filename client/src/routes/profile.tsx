import { createFileRoute } from "@tanstack/react-router";
import Profile from "../components/profile";
import { authCheck } from "../utils/authCheck";

export const Route = createFileRoute("/profile")({
  component: Profile,
  beforeLoad: authCheck,
});
