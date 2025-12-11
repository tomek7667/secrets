import { useState, useEffect, FormEvent } from "react";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { api } from "../../api";
import type { Permission, Token } from "../../types";
import { Table } from "../../components/Table";
import { Button } from "../../components/Button";
import { Input } from "../../components/Input";
import { Modal } from "../../components/Modal";

interface PermissionsPanelProps {
	showToast: (message: string, type: "success" | "error" | "info") => void;
}

export function PermissionsPanel({ showToast }: PermissionsPanelProps) {
	const [permissions, setPermissions] = useState<Permission[]>([]);
	const [tokens, setTokens] = useState<Token[]>([]);
	const [loading, setLoading] = useState(true);

	const [createOpen, setCreateOpen] = useState(false);
	const [createTokenId, setCreateTokenId] = useState("");
	const [createPattern, setCreatePattern] = useState("");
	const [createLoading, setCreateLoading] = useState(false);

	const [editOpen, setEditOpen] = useState(false);
	const [editId, setEditId] = useState("");
	const [editPattern, setEditPattern] = useState("");
	const [editLoading, setEditLoading] = useState(false);

	const load = async () => {
		try {
			const [perms, tkns] = await Promise.all([
				api.permissions.list(),
				api.tokens.list(),
			]);
			setPermissions(perms);
			setTokens(tkns);
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to load permissions",
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
			await api.permissions.create(createTokenId, createPattern);
			showToast("Permission created", "success");
			setCreateOpen(false);
			setCreateTokenId("");
			setCreatePattern("");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to create permission",
				"error"
			);
		} finally {
			setCreateLoading(false);
		}
	};

	const handleEdit = async (e: FormEvent) => {
		e.preventDefault();
		setEditLoading(true);
		try {
			await api.permissions.update(editId, editPattern);
			showToast("Permission updated", "success");
			setEditOpen(false);
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to update permission",
				"error"
			);
		} finally {
			setEditLoading(false);
		}
	};

	const handleDelete = async (id: string) => {
		if (!confirm("Delete this permission?")) return;
		try {
			await api.permissions.delete(id);
			showToast("Permission deleted", "success");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to delete permission",
				"error"
			);
		}
	};

	const openEdit = (perm: Permission) => {
		setEditId(perm.id);
		setEditPattern(perm.secret_key_pattern);
		setEditOpen(true);
	};

	const getTokenPreview = (tokenId: string) => {
		const token = tokens.find((t) => t.id === tokenId);
		if (!token) return tokenId.slice(0, 8) + "...";
		return token.token.slice(0, 12) + "...";
	};

	const columns = [
		{
			key: "token",
			header: "Token",
			render: (p: Permission) => (
				<span className="font-mono text-xs text-slate-400">
					{getTokenPreview(p.token_id)}
				</span>
			),
		},
		{
			key: "pattern",
			header: "Pattern",
			render: (p: Permission) => (
				<span className="font-mono text-sky-400">{p.secret_key_pattern}</span>
			),
		},
		{
			key: "created",
			header: "Created",
			render: (p: Permission) => (
				<span className="text-slate-400">
					{new Date(p.created_at).toLocaleDateString()}
				</span>
			),
		},
		{
			key: "actions",
			header: "",
			className: "text-right w-1",
			render: (p: Permission) => (
				<div className="flex gap-1 justify-end">
					<Button
						variant="ghost"
						size="sm"
						onClick={() => openEdit(p)}
						title="Edit"
					>
						<Pencil size={14} />
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={() => handleDelete(p.id)}
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
				<Button
					onClick={() => setCreateOpen(true)}
					disabled={tokens.length === 0}
				>
					<Plus size={16} />
					New Permission
				</Button>
			</div>

			{loading ? (
				<div className="text-slate-500 py-12 text-center">Loading...</div>
			) : tokens.length === 0 ? (
				<div className="text-slate-500 py-12 text-center">
					Create a token first to add permissions
				</div>
			) : (
				<Table
					columns={columns}
					data={permissions}
					keyField="id"
					emptyMessage="No permissions found"
				/>
			)}

			<Modal
				open={createOpen}
				onClose={() => setCreateOpen(false)}
				title="New Permission"
			>
				<form onSubmit={handleCreate} className="flex flex-col gap-4">
					<div className="flex flex-col gap-1.5">
						<label className="text-xs font-medium text-slate-400 uppercase tracking-wide">
							Token
						</label>
						<select
							value={createTokenId}
							onChange={(e) => setCreateTokenId(e.target.value)}
							required
							className="w-full px-3.5 py-2.5 rounded-lg bg-slate-800 border border-slate-600 text-slate-100 outline-none focus:border-sky-500"
						>
							<option value="">Select a token</option>
							{tokens.map((t) => (
								<option key={t.id} value={t.id}>
									{t.token.slice(0, 20)}...
								</option>
							))}
						</select>
					</div>
					<Input
						id="create-pattern"
						label="Secret Key Pattern"
						value={createPattern}
						onChange={(e) => setCreatePattern(e.target.value)}
						placeholder="e.g. prod/* or DATABASE_*"
						required
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

			<Modal
				open={editOpen}
				onClose={() => setEditOpen(false)}
				title="Edit Permission"
			>
				<form onSubmit={handleEdit} className="flex flex-col gap-4">
					<Input
						id="edit-pattern"
						label="Secret Key Pattern"
						value={editPattern}
						onChange={(e) => setEditPattern(e.target.value)}
						placeholder="e.g. prod/* or DATABASE_*"
						required
					/>
					<div className="flex gap-3 mt-2">
						<Button
							variant="secondary"
							type="button"
							onClick={() => setEditOpen(false)}
							className="flex-1"
						>
							Cancel
						</Button>
						<Button type="submit" loading={editLoading} className="flex-1">
							Save
						</Button>
					</div>
				</form>
			</Modal>
		</div>
	);
}
