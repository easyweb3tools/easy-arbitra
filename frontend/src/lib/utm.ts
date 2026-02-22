const UTM_KEYS = ["utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content"] as const;

export type UTMInput = Record<string, string | undefined>;

export function pickUTM(input: UTMInput): URLSearchParams {
  const out = new URLSearchParams();
  for (const key of UTM_KEYS) {
    const value = input[key];
    if (value && value.trim() !== "") {
      out.set(key, value);
    }
  }
  return out;
}

export function appendUTM(path: string, utm: URLSearchParams): string {
  if (!utm.toString()) return path;
  const joiner = path.includes("?") ? "&" : "?";
  return `${path}${joiner}${utm.toString()}`;
}
