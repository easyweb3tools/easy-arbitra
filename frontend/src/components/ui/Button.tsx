import { type ButtonHTMLAttributes, forwardRef } from "react";
import { Loader2 } from "lucide-react";

type Variant = "filled" | "tinted" | "gray" | "plain" | "destructive";
type Size = "large" | "medium" | "small" | "mini";

const variantClasses: Record<Variant, string> = {
  filled:
    "bg-tint-blue text-white shadow-[0_1px_3px_rgba(0,122,255,0.3)] hover:shadow-[0_4px_12px_rgba(0,122,255,0.3)] hover:brightness-110 active:brightness-95 active:scale-[0.97]",
  tinted:
    "bg-tint-blue/[0.12] text-tint-blue hover:bg-tint-blue/[0.18] active:bg-tint-blue/[0.24] active:scale-[0.97]",
  gray:
    "bg-surface-tertiary text-label-primary hover:bg-surface-tertiary/80 active:scale-[0.97]",
  plain:
    "bg-transparent text-tint-blue hover:text-tint-blue/80 active:text-tint-blue/60 !shadow-none !p-0",
  destructive:
    "bg-tint-red text-white shadow-[0_1px_3px_rgba(255,59,48,0.3)] hover:shadow-[0_4px_12px_rgba(255,59,48,0.3)] hover:brightness-110 active:brightness-95 active:scale-[0.97]",
};

const sizeClasses: Record<Size, string> = {
  large:  "h-[50px] px-7 text-headline rounded-xl",
  medium: "h-11 px-5 text-body rounded-[10px]",
  small:  "h-9 px-4 text-subheadline rounded-lg",
  mini:   "h-7 px-3 text-caption-1 rounded-md",
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
          "transition-all duration-250 ease-apple",
          "disabled:opacity-35 disabled:pointer-events-none disabled:shadow-none",
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
