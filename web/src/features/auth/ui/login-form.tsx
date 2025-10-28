import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { Label } from "@/components/ui/label.tsx";
import { LoginSchema } from "@/features/auth/lib/schemas.ts";
import { cn } from "@/lib/utils.ts";

import { useForm } from "@tanstack/react-form";
import { Link, useNavigate } from "@tanstack/react-router";
import type { LoginCredentials } from "@/types/auth.types.ts";
import { loginFn } from "@/server/auth.tsx";

import { FormInput } from "@/components/form-components/form-input.tsx";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const navigate = useNavigate();
  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
    },

    validators: {
      onSubmit: LoginSchema,
      onSubmitAsync: async ({ value }) => {
        console.log(value);
        try {
          // clearFormErrors(form);
          const loginPayload: LoginCredentials = {
            email: value.username.trim(),
            password: value.password,
          };
          const loggedUser = await loginFn({
            data: {
              email: loginPayload.email,
              password: loginPayload.password,
            },
          });
          console.log(loggedUser);
          // toast.success(
          //   `Welcome back, ${loggedUser.data.resource.profile.full_name || "User"}!`,
          // );
          return undefined;
        } catch (err) {
          // const apiError = handleApiError(err);
          // if (apiError?.errors) {
          //   // handleApiFormErrors(form, apiError);
          // }
          // return apiError?.message || "An unexpected error has occurred.";
        }
      },
    },
    onSubmit: async () => {
      await navigate({ to: "/dashboard" });
    },
  });

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            Sign in with your credentials to continue.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <div className="grid gap-6">
              <div className="grid gap-6">
                <form.Field name={"username"}>
                  {(fieldApi) => (
                    <FormInput
                      fieldApi={fieldApi}
                      label={"Email"}
                      type="email"
                      placeholder="m@example.com"
                    />
                  )}
                </form.Field>

                <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Password</Label>
                    <Link
                      to={"/"}
                      className="text-sm underline-offset-4 hover:underline ml-auto"
                    >
                      Forgot your password?
                    </Link>
                  </div>
                  <form.Field name={"password"}>
                    {(fieldApi) => (
                      <FormInput fieldApi={fieldApi} type="password" />
                    )}
                  </form.Field>
                </div>
                <form.Subscribe
                  selector={(state) => [state.canSubmit, state.isSubmitting]}
                  // biome-ignore lint/correctness/noChildrenProp: <explanation>
                  children={([canSubmit, isSubmitting]) => (
                    <Button
                      type="submit"
                      className="w-full"
                      disabled={!canSubmit || isSubmitting}
                    >
                      Sign In
                    </Button>
                  )}
                />
              </div>
              <div className="flex items-center justify-center text-center text-sm gap-2">
                <span>Don&apos;t have an account?</span>
                <Link
                  to={"/"}
                  className="text-sm underline-offset-4 hover:underline"
                >
                  Sign up
                </Link>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
