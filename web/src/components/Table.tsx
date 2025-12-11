import { ReactNode } from "react";

interface Column<T> {
	key: string;
	header: string;
	render: (item: T) => ReactNode;
	className?: string;
}

interface TableProps<T> {
	columns: Column<T>[];
	data: T[];
	keyField: keyof T;
	emptyMessage?: string;
}

export function Table<T>({
	columns,
	data,
	keyField,
	emptyMessage = "No data",
}: TableProps<T>) {
	return (
		<div className="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50">
			<table className="w-full">
				<thead>
					<tr className="border-b border-slate-700">
						{columns.map((col) => (
							<th
								key={col.key}
								className={`text-left text-xs font-medium text-slate-400 uppercase tracking-wide px-4 py-3 ${
									col.className || ""
								}`}
							>
								{col.header}
							</th>
						))}
					</tr>
				</thead>
				<tbody>
					{data.length === 0 ? (
						<tr>
							<td
								colSpan={columns.length}
								className="px-4 py-12 text-center text-sm text-slate-500"
							>
								{emptyMessage}
							</td>
						</tr>
					) : (
						data.map((item) => (
							<tr
								key={String(item[keyField])}
								className="border-b border-slate-700/50 last:border-0 hover:bg-slate-700/30 transition-colors"
							>
								{columns.map((col) => (
									<td
										key={col.key}
										className={`px-4 py-3 text-sm ${col.className || ""}`}
									>
										{col.render(item)}
									</td>
								))}
							</tr>
						))
					)}
				</tbody>
			</table>
		</div>
	);
}
