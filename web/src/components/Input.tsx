import { InputHTMLAttributes, forwardRef } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
	label?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
	({ label, className = "", id, ...props }, ref) => {
		return (
			<div className="flex flex-col gap-1.5">
				{label && (
					<label
						htmlFor={id}
						className="text-xs font-medium text-slate-400 uppercase tracking-wide"
					>
						{label}
					</label>
				)}
				<input
					ref={ref}
					id={id}
					className={`
            w-full px-3.5 py-2.5 rounded-lg
            bg-slate-800 border border-slate-600 text-slate-100
            placeholder:text-slate-500
            outline-none transition-all duration-150
            focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20
            disabled:opacity-50 disabled:cursor-not-allowed
            ${className}
          `}
					{...props}
				/>
			</div>
		);
	}
);

Input.displayName = "Input";
