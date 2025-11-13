import * as vscode from 'vscode';
import { CipherClient, CallResult } from './cipherClient';
import { parseMetadataAndStrip } from './metadata';
import * as UI from './ui';

// Activate extension: register commands and a simple provider skeleton
export function activate(context: vscode.ExtensionContext) {
  const output = vscode.window.createOutputChannel('Cipher LM Provider');
  context.subscriptions.push(output);

  // Create a CipherClient with defaults (reads settings)
  const config = vscode.workspace.getConfiguration('vscode-lm-provider-cipher');
  const cipherUrl = config.get<string>('cipherAggregatorUrl', 'http://localhost:3020/sse');
  const sessionId = config.get<string>('sessionId') || `session-${Date.now()}-${Math.floor(Math.random()*10000)}`;
  const client = new CipherClient({ url: cipherUrl, sessionId, output });

  // Command: Activate (no-op useful for activationEvents)
  context.subscriptions.push(
    vscode.commands.registerCommand('vscode-lm-provider-cipher.activate', async () => {
      output.appendLine('Cipher LM Provider activated');
      vscode.window.showInformationMessage('Cipher LM Provider activated');
    })
  );

  // Command: Model info (calls list_available_models or tools/list)
  context.subscriptions.push(
    vscode.commands.registerCommand('vscode-lm-provider-cipher.modelInfo', async () => {
      await UI.showModelInfo(client, output);
    })
  );

  // Command: Generate (simple API to test request shaping + metadata parsing)
  context.subscriptions.push(
    vscode.commands.registerCommand('vscode-lm-provider-cipher.generate', async () => {
      const prompt = await vscode.window.showInputBox({ placeHolder: 'Enter a prompt to send to Cipher' });
      if (!prompt) {
        return;
      }

      output.appendLine(`Sending prompt to Cipher: ${prompt}`);

      // Allow per-request tier override via quick pick if enabled in settings
      const allowOverride = config.get<boolean>('allowTierOverride', true);
      const pickedTier = await UI.pickOverrideTier(allowOverride);

      const payload: any = {
        task_description: prompt.split('\n')[0].slice(0, 120),
        messages: [
          { role: 'system', content: 'You are a helpful assistant.' },
          { role: 'user', content: prompt }
        ]
      };

      if (pickedTier && pickedTier !== 'auto') {
        payload.override_tier = pickedTier;
        output.appendLine(`Using override tier: ${pickedTier}`);
      }

      let callResult: CallResult;
      try {
        callResult = await client.callLLMAuto(payload, { stream: false });
      } catch (err: any) {
        output.appendLine(`Cipher call failed: ${String(err)}`);
        vscode.window.showErrorMessage('Cipher call failed: ' + String(err));
        return;
      }

      // Extract metadata and clean response text
      const { text, metadata } = parseMetadataAndStrip(callResult.text);
      output.appendLine('LLM response received.');
      if (metadata) {
        output.appendLine(`Metadata parsed: tier=${metadata.tier} model=${metadata.model} cost=${metadata.cost}`);
      }

      // Show in a new editor tab
      const body = text + (metadata ? `\n\n---\n🤖 [Tier: ${metadata.tier} | Model: ${metadata.model} | Cost: $${metadata.cost} | Tokens: ${metadata.promptTokens}→${metadata.completionTokens} (${metadata.totalTokens})]` : '');
      const doc = await vscode.workspace.openTextDocument({ content: body, language: 'markdown' });
      await vscode.window.showTextDocument(doc, { preview: false });
    })
  );

  // Export client on context for tests or other commands
  (context as any).cipherClient = client;
  output.appendLine('Cipher LM Provider setup complete');
}

// Deactivate
export function deactivate() {
  // nothing for now
}