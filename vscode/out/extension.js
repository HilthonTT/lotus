"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const path = __importStar(require("path"));
const node_1 = require("vscode-languageclient/node");
let client;
async function activate(context) {
    let serverOptions;
    if (process.platform === "win32") {
        // On Windows: launch the Linux binary inside WSL
        // VS Code talks to it over stdio — WSL bridges it transparently
        const binPath = context.asAbsolutePath(path.join("bin", "lotus-lsp-linux"));
        // Convert Windows path to a WSL-accessible path
        const wslPath = binPath
            .replace(/\\/g, "/")
            .replace(/^([A-Z]):/, (_, d) => `/mnt/${d.toLowerCase()}`);
        serverOptions = {
            command: "wsl",
            args: [wslPath],
            transport: node_1.TransportKind.stdio,
        };
    }
    else {
        // On Linux/Mac: run natively
        const serverBin = context.asAbsolutePath(path.join("bin", "lotus-lsp-linux"));
        serverOptions = {
            command: serverBin,
            transport: node_1.TransportKind.stdio,
        };
    }
    const clientOptions = {
        documentSelector: [{ scheme: "file", language: "lotus" }],
    };
    client = new node_1.LanguageClient("lotus-ls", "Lotus Language Server", serverOptions, clientOptions);
    client.setTrace(node_1.Trace.Verbose);
    client.outputChannel.show();
    await client.start();
}
function deactivate() {
    return client?.stop();
}
//# sourceMappingURL=extension.js.map