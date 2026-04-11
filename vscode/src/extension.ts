import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

export async function activate(
  context: vscode.ExtensionContext,
): Promise<void> {
  const out = vscode.window.createOutputChannel("Lotus LSP Debug");
  out.show();
  out.appendLine("Lotus extension activating...");

  const binName = "lotus-lsp-linux";
  const binPath = context.asAbsolutePath(path.join("bin", binName));
  out.appendLine(`Binary path (Windows): ${binPath}`);

  // Convert C:\foo\bar → /mnt/c/foo/bar
  const wslPath = binPath
    .replace(/\\/g, "/")
    .replace(/^([A-Za-z]):/, (_, d) => `/mnt/${d.toLowerCase()}`);
  out.appendLine(`Binary path (WSL): ${wslPath}`);

  const serverOptions: ServerOptions = {
    command: "wsl",
    args: ["bash", "-lc", `chmod +x "${wslPath}" && exec "${wslPath}"`],
    transport: TransportKind.stdio,
    options: {
      env: { ...process.env },
    },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "lotus" }],
    outputChannelName: "Lotus Language Server",
  };

  client = new LanguageClient(
    "lotus-ls",
    "Lotus Language Server",
    serverOptions,
    clientOptions,
  );

  try {
    out.appendLine("Starting LSP client...");
    await client.start();
    out.appendLine("LSP client started successfully!");
  } catch (e) {
    out.appendLine(`LSP client failed to start: ${e}`);
    vscode.window.showErrorMessage(`Lotus LSP failed: ${e}`);
  }
}

export function deactivate(): Thenable<void> | undefined {
  return client?.stop();
}
