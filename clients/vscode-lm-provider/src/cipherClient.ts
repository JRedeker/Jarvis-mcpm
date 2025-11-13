export type CallResult = { text: string };

export interface CipherClientOptions {
  url: string;
  sessionId: string;
  output?: { appendLine: (s: string) => void } | typeof console;
  timeoutMs?: number;
  // Optional injection points for testability
  eventSourceCtor?: any;
  fetchImpl?: any;
}

export class CipherClient {
  private url: string;
  private sessionId: string;
  private output: any;
  private timeoutMs: number;
  // Track active EventSource so we can close from outside or on fallback
  private currentEs: any | null = null;
  // injectable constructors / implementations
  private EventSourceCtor: any;
  private fetchImpl: any;

  constructor(opts: CipherClientOptions) {
    this.url = opts.url.replace(/\/+$/, ''); // trim trailing slash
    this.sessionId = opts.sessionId;
    this.output = opts.output || console;
    this.timeoutMs = opts.timeoutMs ?? 60_000;

    // Use globalThis.require to avoid direct 'require' identifier issues in TS
    // Lazy-injectable implementations for easier unit testing.
    // If not provided, require them at runtime (preserves behavior in production).
    this.EventSourceCtor = opts.eventSourceCtor ?? (globalThis as any).require('eventsource');
    this.fetchImpl = opts.fetchImpl ?? (globalThis as any).require('node-fetch');
  }

  /**
   * Close any active EventSource connection immediately.
   * Useful for test teardown or when switching transport.
   */
  public close() {
    try {
      if (this.currentEs && typeof this.currentEs.close === 'function') {
        this.currentEs.close();
      }
    } catch (e) {
      // ignore
    } finally {
      this.currentEs = null;
    }
  }

  private log(...args: any[]) {
    try {
      if (this.output && typeof this.output.appendLine === 'function') {
        this.output.appendLine(String(args.join(' ')));
      } else {
        console.log(...args);
      }
    } catch {
      // ignore
    }
  }

  // Primary: non-streaming JSON-RPC POST to Cipher Aggregator
  async postJsonRpc(method: string, params: any): Promise<any> {
    const body = {
      jsonrpc: '2.0',
      method,
      id: Date.now(),
      params: Object.assign({ sessionId: this.sessionId }, params || {})
    };

    this.log('POST json-rpc to', this.url, 'body:', JSON.stringify(body).slice(0, 1000));
    const res = await this.fetchImpl(this.url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    });

    if (!res.ok) {
      const text = await res.text().catch(() => '');
      throw new Error(`HTTP ${res.status} ${res.statusText}: ${text}`);
    }

    const json = await res.json().catch(async () => {
      const txt = await res.text();
      throw new Error('Invalid JSON response from Cipher: ' + txt);
    });

    return json;
  }

  // Call llm_inference_auto - non-streaming by default; if opts.stream true, attempt SSE path
  async callLLMAuto(args: any, opts?: { stream?: boolean }): Promise<CallResult> {
    if (opts?.stream) {
      // Attempt SSE path but keep a timeout fallback to POST if SSE fails to produce final result
      try {
        return await this.callLLMAutoSSE(args);
      } catch (err) {
        this.log('SSE failed, falling back to POST:', String(err));
        // fall through to POST fallback
      }
    }

    // POST path (simple): tools/call with name llm_inference_auto
    const rpc = await this.postJsonRpc('tools/call', {
      name: 'llm_inference_auto',
      arguments: args
    });

    // Accept several possible shapes of response
    const result = rpc?.result ?? rpc?.result?.value ?? rpc;
    // Try common locations for text
    const text =
      result?.text ||
      (typeof result === 'string' ? result : null) ||
      (result?.body && (typeof result.body === 'string' ? result.body : result.body.text)) ||
      JSON.stringify(result);

    return { text: String(text ?? '') };
  }

  // SSE entry with simple retry/backoff wrapper
  async callLLMAutoSSE(args: any): Promise<CallResult> {
    const maxAttempts = 3;
    let attempt = 0;
    const backoff = (n: number) => Math.min(1000 * 2 ** n, 8000);

    while (attempt < maxAttempts) {
      attempt++;
      try {
        return await this._callLLMAutoSSEOnce(args);
      } catch (err) {
        this.log(`SSE attempt ${attempt} failed: ${String(err)}`);
        if (attempt >= maxAttempts) throw err;
        await new Promise(r => setTimeout(r, backoff(attempt)));
      }
    }
    throw new Error('SSE failed after retries');
  }

  private async _callLLMAutoSSEOnce(args: any): Promise<CallResult> {
    return new Promise<CallResult>((resolve, reject) => {
      const esUrl = `${this.url}?sessionId=${encodeURIComponent(this.sessionId)}&name=llm_inference_auto`;
      this.log('Opening EventSource to', esUrl);

      let es: any = null;
      let accumulated = '';
      let lastActivity = Date.now();
      const timeoutMs = this.timeoutMs;

      try {
        es = new this.EventSourceCtor(esUrl);
        // remember active ES so callers can close if needed
        this.currentEs = es;
      } catch (err) {
        return reject(new Error('Failed to create EventSource: ' + String(err)));
      }

      const interval = setInterval(() => {
        if (Date.now() - lastActivity > timeoutMs) {
          clearInterval(interval);
          try { if (es && typeof es.close === 'function') es.close(); } catch {}
          this.currentEs = null;
          return reject(new Error('SSE timeout waiting for data'));
        }
      }, Math.min(1000, Math.max(200, Math.floor(timeoutMs / 10))));

      es.onopen = () => {
        lastActivity = Date.now();
        this.log('EventSource opened');
      };

      es.onmessage = (ev: any) => {
        lastActivity = Date.now();
        try {
          const raw = ev.data;
          let obj: any = null;
          try {
            obj = JSON.parse(raw);
          } catch {
            accumulated += raw;
            return;
          }

          // Heartbeat: update activity and ignore payload
          if (obj.event === 'heartbeat' || obj.type === 'heartbeat') {
            lastActivity = Date.now();
            return;
          }

          // Partial/chunk event shapes
          const chunk = obj.text ?? obj.chunk ?? obj.data ?? obj.partial ?? null;
          if (chunk) {
            accumulated += String(chunk);
            // keep accumulating; do not resolve yet
            return;
          }

          // Done/final markers - several possible conventions
          if (obj.event === 'done' || obj.done || obj.type === 'done' || obj.final || obj.complete) {
            accumulated += String(obj.text ?? obj.result ?? '');
            clearInterval(interval);
            try { if (es && typeof es.close === 'function') es.close(); } catch {}
            this.currentEs = null;
            return resolve({ text: accumulated });
          }

          // JSON-RPC final envelope
          if (obj.jsonrpc === '2.0' && obj.result) {
            const res = obj.result;
            const txt = res?.text ?? (typeof res === 'string' ? res : (res?.body?.text ?? JSON.stringify(res)));
            accumulated += String(txt ?? '');
            clearInterval(interval);
            try { if (es && typeof es.close === 'function') es.close(); } catch {}
            this.currentEs = null;
            return resolve({ text: accumulated });
          }

          if (typeof raw === 'string') {
            accumulated += raw;
          }
        } catch (err) {
          clearInterval(interval);
          try { if (es && typeof es.close === 'function') es.close(); } catch {}
          this.currentEs = null;
          return reject(err);
        }
      };

      es.onerror = (err: any) => {
        clearInterval(interval);
        try { if (es && typeof es.close === 'function') es.close(); } catch {}
        this.currentEs = null;
        // Provide a descriptive error so caller can fallback to POST
        return reject(new Error('EventSource error: ' + (err?.message || String(err))));
      };

      // Kick off via POST so aggregator begins streaming to this session
      (async () => {
        try {
          await this.postJsonRpc('tools/call', { name: 'llm_inference_auto', arguments: args });
        } catch (err) {
          clearInterval(interval);
          try { if (es && typeof es.close === 'function') es.close(); } catch {}
          this.currentEs = null;
          reject(err);
        }
      })();
    });
  }

  // Convenience: call the inventory tool to list available models/tiers
  async listAvailableModels(): Promise<string> {
    const rpc = await this.postJsonRpc('tools/call', { name: 'list_available_models', arguments: {} });
    const result = rpc?.result ?? rpc;
    if (typeof result === 'string') return result;
    if (result?.text) return String(result.text);
    return JSON.stringify(result, null, 2);
  }
}

export default CipherClient;