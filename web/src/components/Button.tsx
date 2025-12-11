import { ButtonHTMLAttributes, forwardRef } from "react";
import { Loader2 } from "lucide-react";

type Variant = "primary" | "secondary" | "danger" | "ghost";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
	variant?: Variant;
	loading?: boolean;
	size?: "sm" | "md";
}

const variantClasses: Record<Variant, string> = {
	primary: "bg-sky-500 hover:bg-sky-400 text-white shadow-lg shadow-sky-500/20",
	secondary:
		"bg-slate-700 hover:bg-slate-600 text-slate-200 border border-slate-600",
	danger:
		"bg-red-500/90 hover:bg-red-500 text-white shadow-lg shadow-red-500/20",
	ghost:
		"bg-transparent hover:bg-slate-700/50 text-slate-400 hover:text-slate-200",
};

const sizeClasses: Record<"sm" | "md", string> = {
	sm: "px-2.5 py-1.5 text-xs gap-1.5",
	md: "px-4 py-2.5 text-sm gap-2",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
	(
		{
			variant = "primary",
			size = "md",
			loading = false,
			className = "",
			disabled,
			children,
			...props
		},
		ref
	) => {
		return (
			<button
				ref={ref}
				disabled={disabled || loading}
				className={`
          inline-flex items-center justify-center rounded-lg font-medium
          transition-all duration-150 cursor-pointer
          disabled:opacity-50 disabled:cursor-not-allowed
          ${variantClasses[variant]}
          ${sizeClasses[size]}
          ${className}
        `}
				{...props}
			>
				{loading && (
					<Loader2 size={size === "sm" ? 12 : 16} className="animate-spin" />
				)}
				{children}
			</button>
		);
	}
);

Button.displayName = "Button";
