import { useState, useEffect, useRef, FormEvent } from "react";
import { KeyRound, LogIn } from "lucide-react";
import { Button } from "../components/Button";
import { Input } from "../components/Input";

declare global {
	interface Window {
		turnstile?: {
			render: (
				container: HTMLElement,
				options: {
					sitekey: string;
					callback: (token: string) => void;
					"error-callback"?: () => void;
					"expired-callback"?: () => void;
					theme?: "light" | "dark" | "auto";
				}
			) => string;
			reset: (widgetId: string) => void;
			remove: (widgetId: string) => void;
		};
	}
}

interface LoginProps {
	onLogin: (
		username: string,
		password: string,
		captchaToken?: string
	) => Promise<void>;
	turnstileSiteKey?: string;
}

export function Login({ onLogin, turnstileSiteKey }: LoginProps) {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [captchaToken, setCaptchaToken] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState("");
	const captchaRef = useRef<HTMLDivElement>(null);
	const widgetIdRef = useRef<string | null>(null);

	useEffect(() => {
		if (!turnstileSiteKey || !captchaRef.current) return;

		const renderCaptcha = () => {
			if (window.turnstile && captchaRef.current && !widgetIdRef.current) {
				widgetIdRef.current = window.turnstile.render(captchaRef.current, {
					sitekey: turnstileSiteKey,
					callback: (token: string) => setCaptchaToken(token),
					"expired-callback": () => setCaptchaToken(""),
					"error-callback": () => setCaptchaToken(""),
					theme: "dark",
				});
			}
		};

		if (window.turnstile) {
			renderCaptcha();
		} else {
			const script = document.createElement("script");
			script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js";
			script.async = true;
			script.onload = renderCaptcha;
			document.head.appendChild(script);
		}

		return () => {
			if (widgetIdRef.current && window.turnstile) {
				window.turnstile.remove(widgetIdRef.current);
				widgetIdRef.current = null;
			}
		};
	}, [turnstileSiteKey]);

	const handleSubmit = async (e: FormEvent) => {
		e.preventDefault();
		setError("");

		if (turnstileSiteKey && !captchaToken) {
			setError("Please complete the captcha");
			return;
		}

		setLoading(true);

		try {
			await onLogin(username, password, captchaToken || undefined);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Login failed");
			if (widgetIdRef.current && window.turnstile) {
				window.turnstile.reset(widgetIdRef.current);
				setCaptchaToken("");
			}
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center p-4">
			<div className="w-full max-w-sm animate-slide-up">
				<div className="text-center mb-8">
					<div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-sky-500/10 border border-sky-500/20 mb-4">
						<KeyRound size={32} className="text-sky-400" />
					</div>
					<h1 className="text-2xl font-bold text-slate-100">Secrets</h1>
					<p className="text-sm text-slate-400 mt-1">
						Sign in to manage your secrets
					</p>
				</div>

				<div className="bg-slate-800/50 backdrop-blur-sm rounded-2xl border border-slate-700 p-6">
					<form onSubmit={handleSubmit} className="flex flex-col gap-4">
						<Input
							id="username"
							label="Username"
							type="text"
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							placeholder="Enter username"
							autoComplete="username"
							required
						/>

						<Input
							id="password"
							label="Password"
							type="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							placeholder="Enter password"
							autoComplete="current-password"
							required
						/>

						{turnstileSiteKey && (
							<div className="flex justify-center">
								<div ref={captchaRef} />
							</div>
						)}

						{error && (
							<div className="px-3 py-2 text-sm text-red-300 bg-red-500/10 border border-red-500/20 rounded-lg">
								{error}
							</div>
						)}

						<Button type="submit" loading={loading} className="w-full mt-2">
							<LogIn size={16} />
							Sign In
						</Button>
					</form>
				</div>
			</div>
		</div>
	);
}
