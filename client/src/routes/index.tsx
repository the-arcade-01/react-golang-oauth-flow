import { createFileRoute } from "@tanstack/react-router";
import Navbar from "../components/navbar";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="max-w-screen-lg mx-auto flex flex-cols gap-5 m-5">
      <Navbar />
    </div>
  );
}
