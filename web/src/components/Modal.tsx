import { ReactNode, useEffect, useCallback } from "react";
import { X } from "lucide-react";

interface ModalProps {
	open: boolean;
	onClose: () => void;
	title: string;
	children: ReactNode;
}

export function Modal({ open, onClose, title, children }: ModalProps) {
	const handleKeyDown = useCallback(
		(e: KeyboardEvent) => {
			if (e.key === "Escape") onClose();
		},
		[onClose]
	);

	useEffect(() => {
		if (open) {
			document.addEventListener("keydown", handleKeyDown);
			document.body.style.overflow = "hidden";
		}
		return () => {
			document.removeEventListener("keydown", handleKeyDown);
			document.body.style.overflow = "";
		};
	}, [open, handleKeyDown]);

	if (!open) return null;

	return (
		<div
			className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-in"
			onClick={onClose}
		>
			<div
				className="relative w-full max-w-md bg-slate-800 rounded-2xl border border-slate-700 shadow-2xl animate-slide-up"
				onClick={(e) => e.stopPropagation()}
			>
				<div className="flex items-center justify-between px-6 py-4 border-b border-slate-700">
					<h3 className="text-lg font-semibold text-slate-100">{title}</h3>
					<button
						onClick={onClose}
						className="p-1 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-700 transition-colors cursor-pointer"
					>
						<X size={18} />
					</button>
				</div>
				<div className="p-6">{children}</div>
			</div>
		</div>
	);
}
