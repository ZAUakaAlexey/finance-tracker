import { createFileRoute } from "@tanstack/react-router";
import { createServerFn } from "@tanstack/react-start";
import { authApi } from "@/features/auth/lib/api.ts";
import { useAppSession } from "@/utils/session.ts";
import { LoginForm } from "@/features/auth/ui/login-form.tsx";

export const loginFn = createServerFn({ method: "POST" })
  .inputValidator((data: { email: string; password: string }) => data)
  .handler(async ({ data }) => {
    // Verify credentials (replace with your auth logic)
    const response = await authApi.login({
      email: data.email,
      password: data.password,
    });
    console.log("response", response);
    const user = response?.data?.resource?.user;

    if (!user) {
      return { error: "Invalid credentials" };
    }

    // biome-ignore lint/correctness/useHookAtTopLevel: <explanation>
    const session = await useAppSession();
    await session.update({
      userId: user.id.toString(),
      email: user.email,
      token: response.data.resource.token,
    });

    // throw redirect({ to: "/dashboard" });
  });

export const Route = createFileRoute("/_authed")({
  beforeLoad: ({ context }) => {
    if (!context.user) {
      throw new Error("Not authenticated");
    }
  },
  errorComponent: ({ error }) => {
    if (error.message === "Not authenticated") {
      return <LoginForm />;
    }

    throw error;
  },
});
