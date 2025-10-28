import { createServerFn } from "@tanstack/react-start";
import { redirect } from "@tanstack/react-router";
import { useAppSession } from "@/utils/session.ts";
import { authApi } from "@/features/auth/lib/api.ts";
import { usersApi } from "@/features/users/lib/api.ts";

// Login server function
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

    const session = await useAppSession();
    await session.update({
      userId: user.id.toString(),
      email: user.email,
      token: response.data.resource.token,
    });

    // throw redirect({ to: "/dashboard" });
  });

// Logout server function
export const logoutFn = createServerFn({ method: "POST" }).handler(async () => {
  const session = await useAppSession();
  await session.clear();
  throw redirect({ to: "/" });
});

// Get current user
export const getCurrentUserFn = createServerFn({ method: "GET" }).handler(
  async () => {
    const session = await useAppSession();
    const userId = session.data.userId;

    if (!userId) {
      return null;
    }

    return await usersApi.getUserById(userId);
  },
);
