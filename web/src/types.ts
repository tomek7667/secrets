export interface Secret {
	id: string;
	key: string;
	value: string;
	created_at: string;
	updated_at: string;
}

export interface User {
	id: string;
	username: string;
	password: string;
	created_at: string;
	updated_at: string;
}

export interface Token {
	id: string;
	token: string;
	expires_at: string | null;
	created_at: string;
	updated_at: string;
}

export interface Permission {
	id: string;
	token_id: string;
	secret_key_pattern: string;
	created_at: string;
	updated_at: string;
}

export interface ApiResponse<T> {
	message: string;
	data: T;
}

export interface ApiError {
	message: string;
}
