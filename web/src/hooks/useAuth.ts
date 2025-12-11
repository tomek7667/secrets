import { useState, useCallback } from "react";
import { api } from "../api";

export function useAuth() {
	const [isAuthenticated, setIsAuthenticated] = useState(
		() => !!api.getToken()
	);

	const login = useCallback(
		async (username: string, password: string, captchaToken?: string) => {
			const token = await api.login(username, password, captchaToken);
			api.setToken(token);
			setIsAuthenticated(true);
		},
		[]
	);

	const logout = useCallback(() => {
		api.clearToken();
		setIsAuthenticated(false);
	}, []);

	return { isAuthenticated, login, logout };
}
