import { copyFileSync, mkdirSync } from "fs";
import { dirname, join } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const webDir = dirname(__dirname);
const rootDir = dirname(webDir);

const src = join(webDir, "dist", "index.html");
const destDir = join(rootDir, "internal", "secrets", "dist");
const dest = join(destDir, "index.html");

mkdirSync(destDir, { recursive: true });
copyFileSync(src, dest);
console.log("Copied dist/index.html to internal/secrets/dist/index.html");
