export const dynamic = "force-dynamic";

export default function MethodologyPage() {
  return (
    <section className="space-y-4">
      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-3 text-lg font-semibold">Three-Layer Framework</h2>
        <p className="text-sm text-slate-700">Layer1: PnL facts. Layer2: strategy classification. Layer3: information timing edge.</p>
      </article>
      <article className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-xs text-amber-900">
        <p>Scores are probabilistic estimates.</p>
        <p>Classification is not evidence of wrongdoing.</p>
        <p>This tool is for educational and research purposes only.</p>
      </article>
    </section>
  );
}
