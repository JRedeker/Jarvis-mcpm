export type Role = 'system' | 'user' | 'assistant';

export type ChatMessage = {
  role: Role;
  content: string;
};

export type LLMAutoArgs = {
  task_description: string;
  messages: ChatMessage[];
  override_tier?: string;
  // allow additional pass-through fields if needed
  [k: string]: any;
};

/**
 * Make a concise task description from messages.
 * Strategy: prefer first user message trimmed to 120 chars, fallback to join of first N words.
 */
export function makeTaskDescription(messages: ChatMessage[]): string {
  if (!messages || messages.length === 0) return 'assistant request';
  const firstUser = messages.find(m => m.role === 'user') || messages[0];
  const text = String(firstUser.content || '').trim().replace(/\s+/g, ' ');
  if (text.length <= 120) return text;
  return text.slice(0, 117).trim() + '...';
}

/**
 * Shape a request into the JSON-RPC arguments expected by llm_inference_auto.
 * - messages are normalized to role/content only
 * - task_description is generated if not provided
 */
export function shapeLLMAutoArgs(input: {
  messages?: ChatMessage[] | string;
  task_description?: string;
  override_tier?: string;
}): LLMAutoArgs {
  let messages: ChatMessage[] = [];

  if (typeof input.messages === 'string') {
    messages = [{ role: 'user', content: input.messages }];
  } else if (Array.isArray(input.messages)) {
    messages = input.messages.map(m => ({ role: m.role, content: String(m.content) }));
  } else {
    messages = [{ role: 'user', content: '' }];
  }

  const task_description = input.task_description && input.task_description.trim()
    ? input.task_description.trim()
    : makeTaskDescription(messages);

  const args: LLMAutoArgs = {
    task_description,
    messages
  };

  if (input.override_tier) {
    args.override_tier = input.override_tier;
  }

  return args;
}