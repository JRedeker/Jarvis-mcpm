import { strict as assert } from 'assert';
import { parseMetadataAndStrip } from '../metadata';

describe('metadata.parseMetadataAndStrip', () => {
  it('parses metadata comment and strips it from the response', () => {
    const original = `Here is the assistant output.
<!-- METADATA: tier=m2|model=minimax-01|cost=0.001827|tokens=731,1009,1740 -->`;
    const { text, metadata } = parseMetadataAndStrip(original);

    assert.equal(text.includes('assistant output'), true, 'cleaned text should contain original content');
    assert.ok(metadata, 'metadata should be present');
    assert.equal(metadata?.tier, 'm2');
    assert.equal(metadata?.model, 'minimax-01');
    assert.equal(metadata?.cost, '0.001827');
    assert.equal(metadata?.promptTokens, 731);
    assert.equal(metadata?.completionTokens, 1009);
    assert.equal(metadata?.totalTokens, 1740);
  });

  it('returns original text when no metadata comment present', () => {
    const original = 'No metadata here. Just plain text.';
    const { text, metadata } = parseMetadataAndStrip(original);

    assert.equal(text, original);
    assert.equal(metadata, undefined);
  });

  it('handles malformed or partial metadata gracefully', () => {
    const original = `Output text
<!-- METADATA: tier=m1|model=fast-model|cost=0.0005|tokens=100,abc, -->`;
    const { text, metadata } = parseMetadataAndStrip(original);

    assert.equal(text.includes('Output text'), true);
    assert.ok(metadata, 'metadata parsed even if totals malformed');
    assert.equal(metadata?.tier, 'm1');
    assert.equal(metadata?.model, 'fast-model');
    assert.equal(metadata?.cost, '0.0005');
    assert.equal(metadata?.promptTokens, 100);
    assert.equal(typeof metadata?.completionTokens, 'number');
  });
});