import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

function toWslPath(windowsPath: string): string {
  // C:\Users\foo\bar  →  /mnt/c/Users/foo/bar
  return windowsPath
    .replace(/\\/g, "/")
    .replace(/^([A-Za-z]):/, (_, d) => `/mnt/${d.toLowerCase()}`);
}

export async function activate(
  context: vscode.ExtensionContext,
): Promise<void> {
  const binName = "lotus-lsp-linux";
  const binPath = context.asAbsolutePath(path.join("bin", binName));

  let serverOptions: ServerOptions;

  if (process.platform === "win32") {
    const wslPath = toWslPath(binPath);
    // Show the path so you can verify it in the output panel
    vscode.window.showInformationMessage(`Lotus LSP: launching ${wslPath}`);

    serverOptions = {
      command: "wsl",
      args: ["bash", "-c", `chmod +x "${wslPath}" && "${wslPath}"`],
      transport: TransportKind.stdio,
    };
  } else {
    serverOptions = {
      command: binPath,
      transport: TransportKind.stdio,
    };
  }

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

  await client.start();
}

export function deactivate(): Thenable<void> | undefined {
  return client?.stop();
}
