import { useState, useEffect, FormEvent } from "react";
import { Plus, Copy, Trash2, RefreshCw } from "lucide-react";
import { api } from "../../api";
import type { Token } from "../../types";
import { Table } from "../../components/Table";
import { Button } from "../../components/Button";
import { Input } from "../../components/Input";
import { Modal } from "../../components/Modal";
import { Spoiler } from "../../components/Spoiler";

interface TokensPanelProps {
	showToast: (message: string, type: "success" | "error" | "info") => void;
}

function generateToken(): string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
	const array = new Uint8Array(32);
	crypto.getRandomValues(array);
	return Array.from(array, (b) => chars[b % chars.length]).join("");
}

export function TokensPanel({ showToast }: TokensPanelProps) {
	const [tokens, setTokens] = useState<Token[]>([]);
	const [loading, setLoading] = useState(true);

	const [createOpen, setCreateOpen] = useState(false);
	const [createToken, setCreateToken] = useState("");
	const [createExpires, setCreateExpires] = useState("");
	const [createLoading, setCreateLoading] = useState(false);

	const load = async () => {
		try {
			const data = await api.tokens.list();
			setTokens(data);
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to load tokens",
				"error"
			);
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		load();
	}, []);

	const handleCreate = async (e: FormEvent) => {
		e.preventDefault();
		setCreateLoading(true);
		try {
			const expiresAt = createExpires
				? new Date(createExpires).toISOString()
				: undefined;
			await api.tokens.create(createToken, expiresAt);
			showToast("Token created", "success");
			setCreateOpen(false);
			setCreateToken("");
			setCreateExpires("");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to create token",
				"error"
			);
		} finally {
			setCreateLoading(false);
		}
	};

	const handleDelete = async (id: string) => {
		if (!confirm("Delete this token?")) return;
		try {
			await api.tokens.delete(id);
			showToast("Token deleted", "success");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to delete token",
				"error"
			);
		}
	};

	const copyText = (text: string, msg: string) => {
		navigator.clipboard.writeText(text).then(() => showToast(msg, "success"));
	};

	const openCreate = () => {
		setCreateToken(generateToken());
		setCreateExpires("");
		setCreateOpen(true);
	};

	const columns = [
		{
			key: "token",
			header: "Token",
			render: (t: Token) => <Spoiler value={t.token} />,
		},
		{
			key: "expires",
			header: "Expires",
			render: (t: Token) => (
				<span className="text-slate-400">
					{t.expires_at ? new Date(t.expires_at).toLocaleString() : "Never"}
				</span>
			),
		},
		{
			key: "created",
			header: "Created",
			render: (t: Token) => (
				<span className="text-slate-400">
					{new Date(t.created_at).toLocaleDateString()}
				</span>
			),
		},
		{
			key: "actions",
			header: "",
			className: "text-right w-1",
			render: (t: Token) => (
				<div className="flex gap-1 justify-end">
					<Button
						variant="ghost"
						size="sm"
						onClick={() => copyText(t.token, "Token copied")}
						title="Copy"
					>
						<Copy size={14} />
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={() => handleDelete(t.id)}
						className="text-red-400 hover:text-red-300"
					>
						<Trash2 size={14} />
					</Button>
				</div>
			),
		},
	];

	return (
		<div>
			<div className="flex items-center justify-end mb-4">
				<Button onClick={openCreate}>
					<Plus size={16} />
					New Token
				</Button>
			</div>

			{loading ? (
				<div className="text-slate-500 py-12 text-center">Loading...</div>
			) : (
				<Table
					columns={columns}
					data={tokens}
					keyField="id"
					emptyMessage="No tokens found"
				/>
			)}

			<Modal
				open={createOpen}
				onClose={() => setCreateOpen(false)}
				title="New Token"
			>
				<form onSubmit={handleCreate} className="flex flex-col gap-4">
					<div>
						<Input
							id="create-token"
							label="Token"
							value={createToken}
							onChange={(e) => setCreateToken(e.target.value)}
							placeholder="API Token"
							required
						/>
						<button
							type="button"
							onClick={() => setCreateToken(generateToken())}
							className="inline-flex items-center gap-1 text-xs text-sky-400 hover:text-sky-300 mt-2 cursor-pointer"
						>
							<RefreshCw size={12} />
							Generate new
						</button>
					</div>
					<Input
						id="create-expires"
						label="Expires At (optional)"
						type="datetime-local"
						value={createExpires}
						onChange={(e) => setCreateExpires(e.target.value)}
					/>
					<div className="flex gap-3 mt-2">
						<Button
							variant="secondary"
							type="button"
							onClick={() => setCreateOpen(false)}
							className="flex-1"
						>
							Cancel
						</Button>
						<Button type="submit" loading={createLoading} className="flex-1">
							Create
						</Button>
					</div>
				</form>
			</Modal>
		</div>
	);
}
