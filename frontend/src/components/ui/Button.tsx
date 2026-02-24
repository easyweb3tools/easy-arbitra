import { type ButtonHTMLAttributes, forwardRef } from "react";
import { Loader2 } from "lucide-react";

type Variant = "filled" | "tinted" | "gray" | "plain" | "destructive";
type Size = "large" | "medium" | "small" | "mini";

const variantClasses: Record<Variant, string> = {
  filled:
    "bg-tint-blue text-white hover:brightness-110 active:brightness-90 active:scale-[0.98]",
  tinted:
    "bg-tint-blue/[0.12] text-tint-blue hover:bg-tint-blue/[0.15] active:bg-tint-blue/[0.2] active:scale-[0.98]",
  gray:
    "bg-surface-tertiary text-label-primary hover:brightness-95 active:scale-[0.98]",
  plain:
    "bg-transparent text-tint-blue hover:opacity-70 !shadow-none !p-0",
  destructive:
    "bg-tint-red text-white hover:brightness-110 active:brightness-90 active:scale-[0.98]",
};

const sizeClasses: Record<Size, string> = {
  large:  "h-[50px] px-6 text-headline rounded-md",
  medium: "h-11 px-5 text-body rounded-md",
  small:  "h-9 px-4 text-subheadline rounded-md",
  mini:   "h-7 px-3 text-caption-1 rounded-sm",
};

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: Size;
  loading?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = "filled", size = "medium", loading, children, disabled, className = "", ...props }, ref) => {
    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={[
          "inline-flex items-center justify-center gap-2 font-semibold",
          "transition-all duration-200 ease-apple",
          "disabled:opacity-35 disabled:pointer-events-none",
          "select-none whitespace-nowrap",
          variantClasses[variant],
          sizeClasses[size],
          className,
        ].join(" ")}
        {...props}
      >
        {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
