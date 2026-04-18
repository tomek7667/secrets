import {
	InputHTMLAttributes,
	TextareaHTMLAttributes,
	forwardRef,
} from "react";

type InputProps = (
	| (InputHTMLAttributes<HTMLInputElement> & { multiline?: false })
	| (TextareaHTMLAttributes<HTMLTextAreaElement> & { multiline: true })
) & {
	label?: string;
	id?: string;
};

const sharedClassName = (className: string) => `
            w-full px-3.5 py-2.5 rounded-lg
            bg-slate-800 border border-slate-600 text-slate-100
            placeholder:text-slate-500
            outline-none transition-all duration-150
            focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20
            disabled:opacity-50 disabled:cursor-not-allowed
            ${className}
          `;

export const Input = forwardRef<
	HTMLInputElement | HTMLTextAreaElement,
	InputProps
>(({ label, className = "", id, multiline, ...props }, ref) => {
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
			{multiline ? (
				<textarea
					ref={ref as React.Ref<HTMLTextAreaElement>}
					id={id}
					rows={4}
					className={sharedClassName(className)}
					{...(props as TextareaHTMLAttributes<HTMLTextAreaElement>)}
				/>
			) : (
				<input
					ref={ref as React.Ref<HTMLInputElement>}
					id={id}
					className={sharedClassName(className)}
					{...(props as InputHTMLAttributes<HTMLInputElement>)}
				/>
			)}
		</div>
	);
});

Input.displayName = "Input";
