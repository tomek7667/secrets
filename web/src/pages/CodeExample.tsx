import { useState, type ReactNode } from "react";
import { Check, Copy, FileCode2 } from "lucide-react";
import { Button } from "../components/Button";

const goMainCode = `package main

import (
	"log/slog"

	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/go-multi-logger-slog/logger"
	"github.com/tomek7667/secrets/secretssdk"
)

var secretsClient *secretssdk.Client

func init() {
	logger.SetLogLevel()

	var err error
	secretsClient, err = secretssdk.New(
		utils.Getenv(
			"SECRETS_URL", "http://127.0.0.1:7770",
		),
		utils.Getenv("SECRETS_TOKEN", "s3cr3t-t0k3n"),
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	secret, err := secretsClient.GetSecret("test/key/secret")
	slog.Info(
		"result",
		"secret", secret,
		"err", err,
	)
	allSecrets, err := secretsClient.ListSecrets()
	slog.Info(
		"result",
		"allSecrets", allSecrets,
		"err", err,
	)

	necessaryValueStr := secretsClient.MustGetSecret("test/key/secret").Value
	slog.Info(
		"retrieved secret",
		"str value", necessaryValueStr,
	)
}`;

const goMainCodeLines = goMainCode.split("\n");

const goKeywords = new Set([
	"break",
	"case",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",
]);

const goBuiltins = new Set([
	"append",
	"cap",
	"close",
	"copy",
	"delete",
	"len",
	"make",
	"new",
	"panic",
	"print",
	"println",
	"recover",
]);

const goConstants = new Set(["true", "false", "nil", "iota"]);

const goTokenRegex =
	/\/\/.*$|"(?:\\.|[^"\\])*"|`[^`]*`|\b\d+(?:\.\d+)?\b|\b[A-Za-z_][A-Za-z0-9_]*\b|./g;

const goIdentifierRegex = /^[A-Za-z_][A-Za-z0-9_]*$/;

function getGoTokenClass(
	token: string,
	line: string,
	startIndex: number,
): string {
	if (token.startsWith("//")) return "text-slate-500";
	if (token.startsWith('"') || token.startsWith("`")) return "text-emerald-300";
	if (/^\d/.test(token)) return "text-amber-300";

	if (goKeywords.has(token)) return "text-sky-300";
	if (goBuiltins.has(token)) return "text-cyan-300";
	if (goConstants.has(token)) return "text-orange-300";

	if (goIdentifierRegex.test(token)) {
		const trailing = line.slice(startIndex + token.length).trimStart();
		if (trailing.startsWith("(")) return "text-indigo-300";
		if (/^[A-Z]/.test(token)) return "text-teal-300";
	}

	return "text-slate-200";
}

function highlightGoLine(line: string): ReactNode[] {
	const tokens: ReactNode[] = [];
	goTokenRegex.lastIndex = 0;

	let match = goTokenRegex.exec(line);
	let idx = 0;
	while (match) {
		const token = match[0];
		tokens.push(
			<span
				key={`${match.index}-${idx}`}
				className={getGoTokenClass(token, line, match.index)}
			>
				{token}
			</span>,
		);
		idx += 1;
		match = goTokenRegex.exec(line);
	}

	return tokens;
}

export function CodeExample() {
	const [copied, setCopied] = useState(false);

	const copyCode = async () => {
		try {
			await navigator.clipboard.writeText(goMainCode);
			setCopied(true);
			window.setTimeout(() => setCopied(false), 1600);
		} catch {
			setCopied(false);
		}
	};

	return (
		<section className="mt-4 mb-8 rounded-2xl border border-slate-700/60 bg-slate-900/70 backdrop-blur-sm overflow-hidden animate-slide-up">
			<div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between px-4 py-4 sm:px-5 border-b border-slate-700/60">
				<div className="flex items-center gap-3">
					<div className="h-9 w-9 rounded-xl bg-sky-500/10 border border-sky-500/30 flex items-center justify-center">
						<FileCode2 size={17} className="text-sky-300" />
					</div>
					<div>
						<p className="text-sm font-semibold text-slate-100">
							Golang SDK Quickstart
						</p>
						<p className="text-xs text-slate-400">
							<code>main.go</code> example for fetching a secret securely
						</p>
					</div>
				</div>
				<div className="flex items-center gap-2">
					<span className="px-2.5 py-1 rounded-md text-[11px] font-medium bg-slate-800 text-slate-300 border border-slate-700">
						main.go
					</span>
					<Button
						variant="ghost"
						size="sm"
						onClick={copyCode}
						className="text-slate-300 hover:text-slate-100"
					>
						{copied ? <Check size={14} /> : <Copy size={14} />}
						{copied ? "Copied" : "Copy"}
					</Button>
				</div>
			</div>

			<div className="p-4 sm:p-5">
				<div className="rounded-xl border border-slate-700/80 bg-slate-950/90">
					<pre className="overflow-x-auto p-4 sm:p-5">
						<code className="text-[13px] leading-6">
							{goMainCodeLines.map((line, index) => (
								<div
									key={`${index}-${line}`}
									className="grid grid-cols-[2.25rem_1fr] gap-3"
								>
									<span className="text-right text-slate-500 tabular-nums select-none">
										{index + 1}
									</span>
									<span className="font-mono whitespace-pre">
										{line ? highlightGoLine(line) : " "}
									</span>
								</div>
							))}
						</code>
					</pre>
				</div>
			</div>
		</section>
	);
}
