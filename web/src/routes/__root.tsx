import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  Outlet,
} from "@tanstack/react-router";

import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { TanStackDevtools } from "@tanstack/react-devtools";

// import Header from "../components/Header";

import TanStackQueryDevtools from "../integrations/tanstack-query/devtools";

import StoreDevtools from "../lib/demo-store-devtools";

import appCss from "../styles.css?url";

import type { QueryClient } from "@tanstack/react-query";
import { useAppSession } from "@/utils/session.ts";
// test
import { createServerFn } from "@tanstack/react-start";

interface MyRouterContext {
  queryClient: QueryClient;
}

const fetchUser = createServerFn({ method: "GET" }).handler(async () => {
  // We need to auth on the server so we have access to secure cookies
  const session = await useAppSession();

  if (!session.data.userId) {
    return null;
  }

  return {
    user: session.data,
  };
});

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: async () => {
    const user = await fetchUser();

    return {
      user,
    };
  },
  head: () => ({
    meta: [
      {
        charSet: "utf-8",
      },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      {
        title: "TanStack Start Starter",
      },
    ],
    links: [
      {
        rel: "stylesheet",
        href: appCss,
      },
    ],
  }),

  shellComponent: RootDocument,
});

function RootDocument() {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {/*<Header />*/}
        <Outlet />
        <TanStackDevtools
          config={{
            position: "bottom-right",
          }}
          plugins={[
            {
              name: "Tanstack Router",
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
            StoreDevtools,
          ]}
        />
        <Scripts />
      </body>
    </html>
  );
}
