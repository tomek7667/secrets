import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface SpoilerProps {
	value: string;
	className?: string;
}

export function Spoiler({ value, className = "" }: SpoilerProps) {
	const [visible, setVisible] = useState(false);

	return (
		<button
			onClick={() => setVisible((v) => !v)}
			className={`inline-flex items-center gap-2 text-sm text-slate-300 hover:text-slate-100 transition-colors cursor-pointer ${className}`}
		>
			{visible ? (
				<EyeOff size={14} className="text-slate-500 shrink-0" />
			) : (
				<Eye size={14} className="text-slate-500 shrink-0" />
			)}
			<code className="font-mono w-[180px] text-left overflow-hidden text-ellipsis whitespace-nowrap">
				{visible ? value : "••••••••"}
			</code>
		</button>
	);
}
