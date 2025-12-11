import { KeyRound, Users, Ticket, Shield } from "lucide-react";
import type { Route } from "../hooks/useRouter";

interface Tab {
	id: Route;
	label: string;
	icon: typeof KeyRound;
}

const tabs: Tab[] = [
	{ id: "secrets", label: "Secrets", icon: KeyRound },
	{ id: "users", label: "Users", icon: Users },
	{ id: "tokens", label: "Tokens", icon: Ticket },
	{ id: "permissions", label: "Permissions", icon: Shield },
];

interface TabsProps {
	active: Route;
	onChange: (route: Route) => void;
}

export function Tabs({ active, onChange }: TabsProps) {
	return (
		<div className="flex items-center gap-1 p-1 mb-6 rounded-xl bg-slate-800/50 border border-slate-700/50">
			{tabs.map((tab) => {
				const Icon = tab.icon;
				const isActive = active === tab.id;
				return (
					<button
						key={tab.id}
						onClick={() => onChange(tab.id)}
						className={`
              flex items-center gap-2 flex-1 px-4 py-2.5 rounded-lg text-sm font-medium
              transition-all duration-200 cursor-pointer
              ${
								isActive
									? "bg-sky-500 text-white shadow-lg shadow-sky-500/25"
									: "text-slate-400 hover:text-slate-200 hover:bg-slate-700/50"
							}
            `}
					>
						<Icon size={16} />
						<span className="hidden sm:inline">{tab.label}</span>
					</button>
				);
			})}
		</div>
	);
}
