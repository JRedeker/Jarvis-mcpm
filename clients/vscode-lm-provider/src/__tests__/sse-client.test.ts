import { strict as assert } from 'assert';
import CipherClient from '../cipherClient';

// Minimal mock fetch response helper
function makeFetchOk() {
  return async function fetchImpl(_url: string, _opts: any) {
    return {
      ok: true,
      status: 200,
      statusText: 'OK',
      text: async () => '',
      json: async () => ({ result: {} })
    };
  };
}

describe('cipherClient SSE behavior', () => {
  it('accumulates chunks and resolves on final JSON-RPC envelope', async () => {
    // Mock EventSource that allows sending events programmatically
    class MockEventSource {
      public onopen: (() => void) | null = null;
      public onmessage: ((ev: any) => void) | null = null;
      public onerror: ((err: any) => void) | null = null;
      private closed = false;

      constructor(_url: string) {
        // simulate open shortly after construction
        setTimeout(() => {
          if (this.onopen) this.onopen();
        }, 5);
      }

      sendMessage(obj: any) {
        if (this.onmessage && !this.closed) {
          this.onmessage({ data: JSON.stringify(obj) });
        }
      }

      sendRaw(raw: string) {
        if (this.onmessage && !this.closed) {
          this.onmessage({ data: raw });
        }
      }

      triggerError(err: any) {
        if (this.onerror && !this.closed) {
          this.onerror(err);
        }
      }

      close() {
        this.closed = true;
      }
    }

    const mockFetch = makeFetchOk();
    const client = new CipherClient({
      url: 'http://localhost:3020/sse',
      sessionId: 'session-test-success',
      eventSourceCtor: MockEventSource as any,
      fetchImpl: mockFetch as any,
      timeoutMs: 2000
    });

    // Kick off _callLLMAutoSSEOnce (call the internal single-attempt method)
    const callPromise = (client as any)._callLLMAutoSSEOnce({ task_description: 'test' });

    // Find the active EventSource instance stored on client/currentEs after a tick
    await new Promise(r => setTimeout(r, 20));

    const es = (client as any).currentEs as MockEventSource;
    assert.ok(es, 'EventSource should be created and stored on client.currentEs');

    // Send a partial chunk, then another chunk, then the final JSON-RPC envelope
    es.sendMessage({ text: 'Hello ' });
    es.sendMessage({ text: 'cruel ' });
    // Final JSON-RPC envelope that the client treats as final result
    es.sendMessage({ jsonrpc: '2.0', result: { text: 'world' } });

    const result = await callPromise;
    assert.equal(typeof result.text, 'string');
    assert.equal(result.text.includes('Hello'), true);
    assert.equal(result.text.includes('cruel'), true);
    assert.equal(result.text.includes('world'), true);

    client.close();
  });

  it('rejects on SSE timeout (no messages)', async () => {
    // EventSource that never sends messages (stays silent)
    class SilentEventSource {
      public onopen: (() => void) | null = null;
      public onmessage: ((ev: any) => void) | null = null;
      public onerror: ((err: any) => void) | null = null;
      constructor(_url: string) {
        // open but no messages
        setTimeout(() => { if (this.onopen) this.onopen(); }, 5);
      }
      close() {}
    }

    const client = new CipherClient({
      url: 'http://localhost:3020/sse',
      sessionId: 'session-test-timeout',
      eventSourceCtor: SilentEventSource as any,
      fetchImpl: makeFetchOk() as any,
      timeoutMs: 200 // small timeout for test speed
    });

    try {
      await (client as any)._callLLMAutoSSEOnce({ task_description: 'timeout-test' });
      assert.fail('Expected SSE timeout to reject');
    } catch (err: any) {
      assert.ok(String(err).toLowerCase().includes('timeout'), 'Error should indicate timeout');
    } finally {
      client.close();
    }
  });

  it('rejects when EventSource emits an error', async () => {
    class ErrorEventSource {
      public onopen: (() => void) | null = null;
      public onmessage: ((ev: any) => void) | null = null;
      public onerror: ((err: any) => void) | null = null;
      private closed = false;
      constructor(_url: string) {
        setTimeout(() => { if (this.onopen) this.onopen(); }, 5);
      }
      triggerError(msg: string) {
        if (this.onerror && !this.closed) this.onerror({ message: msg });
      }
      close() { this.closed = true; }
    }

    const mockFetch = makeFetchOk();
    const client = new CipherClient({
      url: 'http://localhost:3020/sse',
      sessionId: 'session-test-error',
      eventSourceCtor: ErrorEventSource as any,
      fetchImpl: mockFetch as any,
      timeoutMs: 2000
    });

    // Start the single-attempt SSE call
    const p = (client as any)._callLLMAutoSSEOnce({ task_description: 'error-test' });

    // Give it a tick, then trigger error on ES
    await new Promise(r => setTimeout(r, 20));
    const es = (client as any).currentEs as ErrorEventSource;
    assert.ok(es, 'EventSource should be created');

    es.triggerError('simulated failure');

    try {
      await p;
      assert.fail('Expected EventSource error to reject the call');
    } catch (err: any) {
      assert.ok(String(err).toLowerCase().includes('eventsource error') || String(err).toLowerCase().includes('simulated failure'));
    } finally {
      client.close();
    }
  });
});