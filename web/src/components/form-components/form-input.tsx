import type { AnyFieldApi } from "@tanstack/react-form";
import * as React from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export interface FormInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  fieldApi?: AnyFieldApi;
  containerClassName?: string;
  showError?: boolean;
}

export const FormInput = React.forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      fieldApi,
      label,
      containerClassName,
      className,
      showError = true,
      ...props
    },
    ref,
  ) => {
    const error = fieldApi?.state?.meta?.errors?.[0]?.toString();

    return (
      <div className={cn("grid gap-3", containerClassName)}>
        {label && (
          <Label
            htmlFor={fieldApi?.name}
            className={cn(error && "text-destructive")}
          >
            {label}
          </Label>
        )}
        <Input
          ref={ref}
          id={fieldApi?.name}
          name={fieldApi?.name}
          value={fieldApi?.state?.value || ""}
          onBlur={fieldApi?.handleBlur}
          onChange={(e) => fieldApi?.handleChange?.(e.target.value)}
          className={cn(
            error && "border-destructive focus-visible:ring-destructive",
            className,
          )}
          aria-invalid={!!error}
          {...props}
        />
        {showError && error && (
          <p className="text-sm text-destructive">{error}</p>
        )}
      </div>
    );
  },
);

FormInput.displayName = "FormInput";
