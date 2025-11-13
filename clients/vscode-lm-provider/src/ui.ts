import * as vscode from 'vscode';
import { CipherClient } from './cipherClient';

/**
 * Show a quick pick for tier override. Returns the selected tier id (e.g., m2) or undefined if cancelled.
 */
export async function pickOverrideTier(allowOverride: boolean): Promise<string | undefined> {
  if (!allowOverride) return undefined;

  const items: vscode.QuickPickItem[] = [
    { label: 'Auto (Cipher)', description: 'Let Cipher pick the tier' },
    { label: 'Tier L0', description: 'l0 (speed)' },
    { label: 'Tier M1', description: 'm1 (light)' },
    { label: 'Tier M2', description: 'm2 (balanced)' },
    { label: 'Tier M3', description: 'm3 (deep reasoning)' },
    { label: 'Tier M4', description: 'm4 (max reasoning)' }
  ];

  const picked = await vscode.window.showQuickPick(items, {
    placeHolder: 'Select a tier override for this request (or pick Auto)',
    ignoreFocusOut: true
  });

  if (!picked) return undefined;
  const map: Record<string, string> = {
    'Auto (Cipher)': 'auto',
    'Tier L0': 'l0',
    'Tier M1': 'm1',
    'Tier M2': 'm2',
    'Tier M3': 'm3',
    'Tier M4': 'm4'
  };
  return map[picked.label] ?? undefined;
}

/**
 * Fetch model info from Cipher via the client and render it in a new editor.
 */
export async function showModelInfo(client: CipherClient, output?: { appendLine: (s: string) => void }) {
  try {
    if (output && typeof output.appendLine === 'function') {
      output.appendLine('Fetching model info from Cipher...');
    } else {
      console.log('Fetching model info from Cipher...');
    }

    const raw = await client.listAvailableModels();
    const title = 'Cipher Model Info';
    const content = typeof raw === 'string' ? raw : JSON.stringify(raw, null, 2);
    const doc = await vscode.workspace.openTextDocument({ content: `# ${title}\n\n${content}`, language: 'markdown' });
    await vscode.window.showTextDocument(doc, { preview: true });

    if (output && typeof output.appendLine === 'function') {
      output.appendLine('Model info displayed');
    } else {
      console.log('Model info displayed');
    }
  } catch (err: any) {
    if (output && typeof output.appendLine === 'function') {
      output.appendLine('Error fetching model info: ' + String(err));
    } else {
      console.error('Error fetching model info: ' + String(err));
    }
    void vscode.window.showErrorMessage('Failed to fetch model info from Cipher: ' + String(err));
  }
}

export default {
  pickOverrideTier,
  showModelInfo
};