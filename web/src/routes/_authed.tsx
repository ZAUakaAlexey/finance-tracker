import { createFileRoute, Outlet } from "@tanstack/react-router";
import { createServerFn } from "@tanstack/react-start";
import { authApi } from "@/features/auth/lib/api.ts";
import { useAppSession } from "@/utils/session.ts";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar.tsx";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb.tsx";

import { Separator } from "@/components/ui/separator";
import { AppSidebar } from "@/components/app-sidebar.tsx";
import type { IUser } from "@/types/auth.types.ts";
import { ModeToggle } from "@/components/mode-toggle.tsx";
import LoginPage from "@/features/auth/pub/login-page.tsx";

export const loginFn = createServerFn({ method: "POST" })
  .inputValidator((data: { email: string; password: string }) => data)
  .handler(async ({ data }) => {
    // Verify credentials (replace with your auth logic)
    const response = await authApi.login({
      email: data.email,
      password: data.password,
    });
    console.log("response", response); //TODO remove in PROD
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
      return <LoginPage />;
    }

    throw error;
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  const { user } = Route.useRouteContext();
  return (
    <SidebarProvider>
      <AppSidebar loggedUser={user?.user as IUser} />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
          <div className="flex flex-1 items-center justify-between pr-4">
            <div className="flex items-center gap-2 px-4">
              <SidebarTrigger className="-ml-1" />
              <Separator
                orientation="vertical"
                className="mr-2 data-[orientation=vertical]:h-4"
              />
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem className="hidden md:block">
                    <BreadcrumbLink href="#">
                      Building Your Application
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator className="hidden md:block" />
                  <BreadcrumbItem>
                    <BreadcrumbPage>Data Fetching</BreadcrumbPage>
                  </BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </div>
            <ModeToggle />
          </div>
        </header>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  );
}
