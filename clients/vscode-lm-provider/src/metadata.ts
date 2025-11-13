export type Metadata = {
  tier: string;
  model: string;
  cost: string;
  promptTokens: number;
  completionTokens: number;
  totalTokens: number;
};

export function parseMetadataAndStrip(original: string): { text: string; metadata?: Metadata } {
  if (!original) return { text: original };

  // Regex tolerant to whitespace and small variations; captures the six fields in order.
  const re = /<!--\s*METADATA:\s*tier=([^|]+)\|model=([^|]+)\|cost=([^|]+)\|tokens=([^,]+),([^,]+),([^>]+)\s*-->/i;
  const match = original.match(re);

  if (!match) {
    return { text: original };
  }

  const [, tierRaw, modelRaw, costRaw, promptRaw, completionRaw, totalRaw] = match;
  const tier = tierRaw.trim();
  const model = modelRaw.trim();
  const cost = costRaw.trim();
  const promptTokens = parseInt(promptRaw.trim(), 10) || 0;
  const completionTokens = parseInt(completionRaw.trim(), 10) || 0;
  // totalRaw may include trailing characters; extract leading number
  const totalTokens = parseInt((totalRaw || '').trim().split(/\D/)[0], 10) || (promptTokens + completionTokens);

  // Strip the metadata comment from the original response
  const cleaned = original.replace(re, '').trim();

  const metadata: Metadata = {
    tier,
    model,
    cost,
    promptTokens,
    completionTokens,
    totalTokens
  };

  return { text: cleaned, metadata };
}