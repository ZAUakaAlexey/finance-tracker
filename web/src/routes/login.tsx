import { createFileRoute } from "@tanstack/react-router";
import { LoginForm } from "@/features/auth/ui/login-form.tsx";

export const Route = createFileRoute("/login")({
  component: RouteComponent,
});

function RouteComponent() {
  return <LoginForm />;
}
