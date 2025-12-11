import type { ApiResponse, Secret, User, Token, Permission } from "./types";

const getToken = (): string | null => localStorage.getItem("jwt");

const setToken = (token: string) => localStorage.setItem("jwt", token);

const clearToken = () => localStorage.removeItem("jwt");

async function request<T>(
	method: string,
	url: string,
	body?: unknown
): Promise<T> {
	const token = getToken();
	const headers: Record<string, string> = {
		"Content-Type": "application/json",
	};
	if (token) {
		headers["Authorization"] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		method,
		headers,
		body: body ? JSON.stringify(body) : undefined,
	});

	const data = await response.json();

	if (!response.ok) {
		throw new Error(data.message || "Request failed");
	}

	return data.data;
}

export const api = {
	getToken,
	setToken,
	clearToken,

	login: async (
		username: string,
		password: string,
		captchaToken?: string
	): Promise<string> => {
		const response = await fetch("/login", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username,
				password,
				captchaToken: captchaToken || "",
			}),
		});
		const data: ApiResponse<{ token: string }> = await response.json();
		if (!response.ok) {
			throw new Error(
				(data as unknown as { message: string }).message || "Login failed"
			);
		}
		return data.data.token;
	},

	secrets: {
		list: () => request<Secret[]>("GET", "/api/secrets"),
		create: (key: string, value: string) =>
			request<Secret>("POST", "/api/secrets", { key, value }),
		update: (key: string, value: string) =>
			request<Secret>("PUT", `/api/secrets?key=${encodeURIComponent(key)}`, {
				value,
			}),
		delete: (key: string) =>
			request<void>("DELETE", `/api/secrets?key=${encodeURIComponent(key)}`),
	},

	users: {
		list: () => request<User[]>("GET", "/api/users"),
		create: (username: string, password: string) =>
			request<User>("POST", "/api/users", { username, password }),
		delete: (id: string) =>
			request<void>("DELETE", `/api/users/${encodeURIComponent(id)}`),
	},

	tokens: {
		list: () => request<Token[]>("GET", "/api/tokens"),
		create: (token: string, expiresAt?: string) =>
			request<Token>("POST", "/api/tokens", {
				token,
				expires_at: expiresAt || null,
			}),
		delete: (id: string) =>
			request<void>("DELETE", `/api/tokens/${encodeURIComponent(id)}`),
	},

	permissions: {
		list: () => request<Permission[]>("GET", "/api/permissions"),
		create: (tokenId: string, secretKeyPattern: string) =>
			request<Permission>("POST", "/api/permissions", {
				token_id: tokenId,
				secret_key_pattern: secretKeyPattern,
			}),
		update: (id: string, secretKeyPattern: string) =>
			request<Permission>("PUT", `/api/permissions/${encodeURIComponent(id)}`, {
				secret_key_pattern: secretKeyPattern,
			}),
		delete: (id: string) =>
			request<void>("DELETE", `/api/permissions/${encodeURIComponent(id)}`),
	},
};
