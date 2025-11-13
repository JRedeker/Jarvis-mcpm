/**
 * Lightweight telemetry helper for the VS Code Cipher LM provider.
 *
 * Responsibilities:
 * - Respect telemetryEnabled flag (do nothing when disabled)
 * - Provide recordEvent(name, props) and timedEvent helpers
 * - Minimal anonymization helper (hashes identifiers) so tests can assert presence
 * - Expose an in-memory sink for unit tests (getEvents) and a flush hook
 */

export type TelemetryEvent = {
  name: string;
  timestamp: number;
  durationMs?: number;
  props?: Record<string, any>;
};

export type TelemetryOptions = {
  enabled?: boolean; // default true
  output?: { appendLine: (s: string) => void } | typeof console;
  // hook called when events are emitted; useful for CI/analytics integration
  onEmit?: (ev: TelemetryEvent) => void;
};

function simpleHash(s: string): string {
  // Very small, deterministic hash for anonymization in tests.
  // Not cryptographically secure; purpose is to avoid storing raw identifiers.
  let h = 0;
  for (let i = 0; i < s.length; i++) {
    h = (h << 5) - h + s.charCodeAt(i);
    h |= 0;
  }
  return 'h' + Math.abs(h).toString(36);
}

export class Telemetry {
  private enabled: boolean;
  private sink: TelemetryEvent[] = [];
  private output: any;
  private onEmit?: (ev: TelemetryEvent) => void;

  constructor(opts?: TelemetryOptions) {
    this.enabled = opts?.enabled ?? true;
    this.output = opts?.output ?? console;
    this.onEmit = opts?.onEmit;
  }

  public isEnabled() {
    return this.enabled;
  }

  public setEnabled(v: boolean) {
    this.enabled = Boolean(v);
  }

  // Record a generic event. Props will be shallow-copied and sanitized.
  public recordEvent(name: string, props?: Record<string, any>) {
    if (!this.enabled) return;
    const ev: TelemetryEvent = {
      name,
      timestamp: Date.now(),
      props: this.sanitizeProps(props)
    };
    this.sink.push(ev);
    try {
      if (this.output && typeof this.output.appendLine === 'function') {
        this.output.appendLine(`TELEMETRY: ${name} ${JSON.stringify(ev.props || {})}`);
      } else {
        // fallback to console
        console.log('TELEMETRY:', name, ev.props || {});
      }
    } catch {
      // swallow output errors
    }
    if (this.onEmit) {
      try { this.onEmit(ev); } catch {}
    }
  }

  // Start a timed event and return an end() function that records duration.
  public startTimedEvent(name: string, props?: Record<string, any>) {
    const start = Date.now();
    const self = this;
    return {
      end(endProps?: Record<string, any>) {
        const durationMs = Date.now() - start;
        const merged = Object.assign({}, props || {}, endProps || {}, { durationMs });
        self.recordEvent(name, merged);
      }
    };
  }

  // Expose events for testing/inspection
  public getEvents(): TelemetryEvent[] {
    return this.sink.slice();
  }

  public clear() {
    this.sink = [];
  }

  // Simple sanitization: remove prompt-like fields and hash identifiers
  private sanitizeProps(props?: Record<string, any>): Record<string, any> | undefined {
    if (!props) return undefined;
    const out: Record<string, any> = {};
    for (const k of Object.keys(props)) {
      const v = props[k];
      if (k.toLowerCase().includes('prompt') || k.toLowerCase().includes('content') || k.toLowerCase().includes('text')) {
        // redact prompt-like content
        out[k] = '[REDACTED]';
      } else if (typeof v === 'string' && (k.toLowerCase().includes('session') || k.toLowerCase().includes('id') || k.toLowerCase().includes('user'))) {
        out[k] = simpleHash(v);
      } else {
        out[k] = v;
      }
    }
    return out;
  }
}

export default Telemetry;