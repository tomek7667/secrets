import { useState, useEffect, FormEvent } from "react";
import {
	Plus,
	Pencil,
	Trash2,
	Search,
	KeyRound,
	ClipboardCopy,
} from "lucide-react";
import { api } from "../../api";
import type { Secret } from "../../types";
import { Table } from "../../components/Table";
import { Button } from "../../components/Button";
import { Input } from "../../components/Input";
import { Modal } from "../../components/Modal";
import { Spoiler } from "../../components/Spoiler";

interface SecretsPanelProps {
	showToast: (message: string, type: "success" | "error" | "info") => void;
}

export function SecretsPanel({ showToast }: SecretsPanelProps) {
	const [secrets, setSecrets] = useState<Secret[]>([]);
	const [loading, setLoading] = useState(true);
	const [filter, setFilter] = useState("");

	const [createOpen, setCreateOpen] = useState(false);
	const [createKey, setCreateKey] = useState("");
	const [createValue, setCreateValue] = useState("");
	const [createLoading, setCreateLoading] = useState(false);

	const [editOpen, setEditOpen] = useState(false);
	const [editKey, setEditKey] = useState("");
	const [editValue, setEditValue] = useState("");
	const [editLoading, setEditLoading] = useState(false);

	const load = async () => {
		try {
			const data = await api.secrets.list();
			setSecrets(data);
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to load secrets",
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
			await api.secrets.create(createKey, createValue);
			showToast("Secret created", "success");
			setCreateOpen(false);
			setCreateKey("");
			setCreateValue("");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to create secret",
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
			await api.secrets.update(editKey, editValue);
			showToast("Secret updated", "success");
			setEditOpen(false);
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to update secret",
				"error"
			);
		} finally {
			setEditLoading(false);
		}
	};

	const handleDelete = async (key: string) => {
		if (!confirm(`Delete secret "${key}"?`)) return;
		try {
			await api.secrets.delete(key);
			showToast("Secret deleted", "success");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to delete secret",
				"error"
			);
		}
	};

	const copyText = (text: string, msg: string) => {
		navigator.clipboard.writeText(text).then(() => showToast(msg, "success"));
	};

	const openEdit = (secret: Secret) => {
		setEditKey(secret.key);
		setEditValue(atob(secret.value));
		setEditOpen(true);
	};

	const filtered = secrets.filter((s) =>
		s.key.toLowerCase().includes(filter.toLowerCase())
	);

	const columns = [
		{
			key: "key",
			header: "Key",
			render: (s: Secret) => (
				<span className="font-mono text-sky-400">{s.key}</span>
			),
		},
		{
			key: "value",
			header: "Value",
			render: (s: Secret) => <Spoiler value={atob(s.value)} />,
		},
		{
			key: "actions",
			header: "",
			className: "text-right w-1",
			render: (s: Secret) => (
				<div className="flex items-center gap-1 justify-end">
					<Button
						variant="ghost"
						size="sm"
						onClick={() => copyText(s.key, "Key copied")}
						title="Copy key"
					>
						<KeyRound size={14} />
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={() => copyText(atob(s.value), "Value copied")}
						title="Copy value"
					>
						<ClipboardCopy size={14} />
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={() => openEdit(s)}
						title="Edit"
					>
						<Pencil size={14} />
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={() => handleDelete(s.key)}
						title="Delete"
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
			<div className="flex items-center justify-between gap-4 mb-4 flex-wrap">
				<div className="relative">
					<Search
						size={16}
						className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500"
					/>
					<input
						type="text"
						placeholder="Filter secrets..."
						value={filter}
						onChange={(e) => setFilter(e.target.value)}
						className="pl-9 pr-4 py-2 rounded-lg bg-slate-800 border border-slate-700 text-sm text-slate-200 placeholder:text-slate-500 outline-none focus:border-sky-500 w-64"
					/>
				</div>
				<Button onClick={() => setCreateOpen(true)}>
					<Plus size={16} />
					New Secret
				</Button>
			</div>

			{loading ? (
				<div className="text-slate-500 py-12 text-center">Loading...</div>
			) : (
				<Table
					columns={columns}
					data={filtered}
					keyField="id"
					emptyMessage="No secrets found"
				/>
			)}

			<Modal
				open={createOpen}
				onClose={() => setCreateOpen(false)}
				title="New Secret"
			>
				<form onSubmit={handleCreate} className="flex flex-col gap-4">
					<Input
						id="create-key"
						label="Key"
						value={createKey}
						onChange={(e) => setCreateKey(e.target.value)}
						placeholder="e.g. DATABASE_URL"
						required
					/>
					<Input
						id="create-value"
						label="Value"
						value={createValue}
						onChange={(e) => setCreateValue(e.target.value)}
						placeholder="Secret value"
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
				title="Edit Secret"
			>
				<form onSubmit={handleEdit} className="flex flex-col gap-4">
					<Input id="edit-key" label="Key" value={editKey} disabled />
					<Input
						id="edit-value"
						label="Value"
						value={editValue}
						onChange={(e) => setEditValue(e.target.value)}
						placeholder="New value"
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
