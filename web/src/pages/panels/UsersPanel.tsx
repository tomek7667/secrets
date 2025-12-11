import { useState, useEffect, FormEvent } from "react";
import { Plus, Trash2 } from "lucide-react";
import { api } from "../../api";
import type { User } from "../../types";
import { Table } from "../../components/Table";
import { Button } from "../../components/Button";
import { Input } from "../../components/Input";
import { Modal } from "../../components/Modal";

interface UsersPanelProps {
	showToast: (message: string, type: "success" | "error" | "info") => void;
}

export function UsersPanel({ showToast }: UsersPanelProps) {
	const [users, setUsers] = useState<User[]>([]);
	const [loading, setLoading] = useState(true);

	const [createOpen, setCreateOpen] = useState(false);
	const [createUsername, setCreateUsername] = useState("");
	const [createPassword, setCreatePassword] = useState("");
	const [createLoading, setCreateLoading] = useState(false);

	const load = async () => {
		try {
			const data = await api.users.list();
			setUsers(data);
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to load users",
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
			await api.users.create(createUsername, createPassword);
			showToast("User created", "success");
			setCreateOpen(false);
			setCreateUsername("");
			setCreatePassword("");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to create user",
				"error"
			);
		} finally {
			setCreateLoading(false);
		}
	};

	const handleDelete = async (id: string) => {
		if (!confirm("Delete this user?")) return;
		try {
			await api.users.delete(id);
			showToast("User deleted", "success");
			load();
		} catch (err) {
			showToast(
				err instanceof Error ? err.message : "Failed to delete user",
				"error"
			);
		}
	};

	const columns = [
		{
			key: "username",
			header: "Username",
			render: (u: User) => (
				<span className="font-medium text-slate-200">{u.username}</span>
			),
		},
		{
			key: "id",
			header: "ID",
			render: (u: User) => (
				<span className="font-mono text-xs text-slate-500">
					{u.id.slice(0, 8)}...
				</span>
			),
		},
		{
			key: "created",
			header: "Created",
			render: (u: User) => (
				<span className="text-slate-400">
					{new Date(u.created_at).toLocaleDateString()}
				</span>
			),
		},
		{
			key: "actions",
			header: "",
			className: "text-right w-1",
			render: (u: User) => (
				<Button
					variant="ghost"
					size="sm"
					onClick={() => handleDelete(u.id)}
					className="text-red-400 hover:text-red-300"
				>
					<Trash2 size={14} />
				</Button>
			),
		},
	];

	return (
		<div>
			<div className="flex items-center justify-end mb-4">
				<Button onClick={() => setCreateOpen(true)}>
					<Plus size={16} />
					New User
				</Button>
			</div>

			{loading ? (
				<div className="text-slate-500 py-12 text-center">Loading...</div>
			) : (
				<Table
					columns={columns}
					data={users}
					keyField="id"
					emptyMessage="No users found"
				/>
			)}

			<Modal
				open={createOpen}
				onClose={() => setCreateOpen(false)}
				title="New User"
			>
				<form onSubmit={handleCreate} className="flex flex-col gap-4">
					<Input
						id="create-username"
						label="Username"
						value={createUsername}
						onChange={(e) => setCreateUsername(e.target.value)}
						placeholder="Username"
						required
					/>
					<Input
						id="create-password"
						label="Password"
						type="password"
						value={createPassword}
						onChange={(e) => setCreatePassword(e.target.value)}
						placeholder="Password"
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
		</div>
	);
}
