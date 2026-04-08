import { KeyRound, LogOut, ArrowDown, ArrowUp } from "lucide-react";
import { Tabs } from "../components/Tabs";
import { Button } from "../components/Button";
import { SecretsPanel } from "./panels/SecretsPanel";
import { UsersPanel } from "./panels/UsersPanel";
import { TokensPanel } from "./panels/TokensPanel";
import { PermissionsPanel } from "./panels/PermissionsPanel";
import type { Route } from "../hooks/useRouter";
import { CodeExample } from "./CodeExample";
import { useState } from "react";

interface DashboardProps {
	route: Route;
	onRouteChange: (route: Route) => void;
	onLogout: () => void;
	showToast: (message: string, type: "success" | "error" | "info") => void;
}

export function Dashboard({
	route,
	onRouteChange,
	onLogout,
	showToast,
}: DashboardProps) {
	const [codeExampleOpened, setCodeExampleOpened] = useState(false);

	return (
		<div className="min-h-screen">
			<div className="max-w-6xl mx-auto p-4 sm:p-6 lg:p-8">
				<header className="flex items-center justify-between mb-8">
					<div className="flex items-center gap-3">
						<div className="flex items-center justify-center w-10 h-10 rounded-xl bg-sky-500/10 border border-sky-500/20">
							<KeyRound size={20} className="text-sky-400" />
						</div>
						<div>
							<h1 className="text-xl font-bold text-slate-100">Secrets</h1>
							<p className="text-xs text-slate-500">
								Manage secrets, users & permissions
							</p>
						</div>
					</div>
					<Button variant="ghost" size="sm" onClick={onLogout}>
						<LogOut size={14} />
						Logout
					</Button>
				</header>

				<Tabs active={route} onChange={onRouteChange} />
				<Button
					variant="secondary"
					size="md"
					onClick={() => setCodeExampleOpened((open) => !open)}
					style={{
						marginBottom: 8,
					}}
				>
					{codeExampleOpened ? "Hide" : "Show"} example SDK code
					{codeExampleOpened ? <ArrowUp /> : <ArrowDown />}
				</Button>

				{codeExampleOpened ? <CodeExample /> : null}

				<main className="animate-fade-in">
					{route === "secrets" && <SecretsPanel showToast={showToast} />}
					{route === "users" && <UsersPanel showToast={showToast} />}
					{route === "tokens" && <TokensPanel showToast={showToast} />}
					{route === "permissions" && (
						<PermissionsPanel showToast={showToast} />
					)}
				</main>
			</div>
		</div>
	);
}
