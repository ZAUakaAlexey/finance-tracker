import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/_authed/dashboard")({
  component: DashboardComponent,
});

function DashboardComponent() {
  const { user } = Route.useRouteContext();
  const navigate = useNavigate();
  return (
    <div>
      <h1>Welcome, {user?.user.email}!</h1>
      {/* Dashboard content */}
      <div className="h-screen flex items-center justify-center">
        <Button onClick={async () => await navigate({ to: "/logout" })}>
          Click Me
        </Button>
      </div>
    </div>
  );
}
