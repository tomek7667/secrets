import { CheckCircle, XCircle, Info } from "lucide-react";
import type {
	Toast as ToastType,
	ToastType as ToastVariant,
} from "../hooks/useToast";

interface ToastContainerProps {
	toasts: ToastType[];
}

const config: Record<
	ToastVariant,
	{ icon: typeof CheckCircle; bg: string; border: string }
> = {
	success: {
		icon: CheckCircle,
		bg: "bg-emerald-500/10",
		border: "border-emerald-500/50",
	},
	error: {
		icon: XCircle,
		bg: "bg-red-500/10",
		border: "border-red-500/50",
	},
	info: {
		icon: Info,
		bg: "bg-sky-500/10",
		border: "border-sky-500/50",
	},
};

export function ToastContainer({ toasts }: ToastContainerProps) {
	return (
		<div className="fixed right-4 bottom-4 z-[1000] flex flex-col gap-2">
			{toasts.map((toast) => {
				const { icon: Icon, bg, border } = config[toast.type];
				return (
					<div
						key={toast.id}
						className={`
              flex items-center gap-3 px-4 py-3 rounded-xl
              text-sm text-slate-200 border backdrop-blur-sm
              animate-slide-up shadow-lg
              ${bg} ${border}
            `}
					>
						<Icon
							size={18}
							className={
								toast.type === "success"
									? "text-emerald-400"
									: toast.type === "error"
									? "text-red-400"
									: "text-sky-400"
							}
						/>
						{toast.message}
					</div>
				);
			})}
		</div>
	);
}
