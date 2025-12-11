import { useAuth } from "./hooks/useAuth";
import { useRouter } from "./hooks/useRouter";
import { useToast } from "./hooks/useToast";
import { Login } from "./pages/Login";
import { Dashboard } from "./pages/Dashboard";
import { ToastContainer } from "./components/Toast";

declare global {
	interface Window {
		TURNSTILE_SITE_KEY?: string;
	}
}

export default function App() {
	const { isAuthenticated, login, logout } = useAuth();
	const { route, setRoute } = useRouter();
	const { toasts, showToast } = useToast();

	const turnstileSiteKey = window.TURNSTILE_SITE_KEY;

	return (
		<>
			{isAuthenticated ? (
				<Dashboard
					route={route}
					onRouteChange={setRoute}
					onLogout={logout}
					showToast={showToast}
				/>
			) : (
				<Login onLogin={login} turnstileSiteKey={turnstileSiteKey} />
			)}
			<ToastContainer toasts={toasts} />
		</>
	);
}
