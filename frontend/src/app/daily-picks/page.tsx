import Link from "next/link";
import { Sparkles, TrendingUp, TrendingDown, Brain, Clock, CheckCircle, Loader2 } from "lucide-react";
import { getDailyPick, getDailyPickHistory, getNovaSessions } from "@/lib/api";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import type { NovaSession } from "@/lib/types";

export const dynamic = "force-dynamic";

function PhaseIcon({ phase }: { phase: string }) {
    if (phase === "final") return <CheckCircle className="h-4 w-4 text-tint-green" />;
    return <Loader2 className="h-4 w-4 text-tint-blue animate-spin" />;
}

function SessionCard({ session, locale }: { session: NovaSession; locale: string }) {
    const obs = session.observations_json as { notes?: string } | null;
    const observationText = obs?.notes || "";

    return (
        <div className={[
            "relative rounded-xl p-4 transition-all",
            session.phase === "final"
                ? "bg-gradient-to-r from-tint-green/10 to-tint-blue/5 border border-tint-green/20"
                : "bg-surface-secondary border border-separator/30",
        ].join(" ")}>
            <div className="flex items-center gap-2 mb-2">
                <PhaseIcon phase={session.phase} />
                <span className="text-caption-1 font-bold text-label-secondary uppercase tracking-wide">
                    Round {session.round}
                </span>
                <span className="text-caption-2 text-label-quaternary">
                    {new Date(session.created_at).toLocaleTimeString()}
                </span>
                {session.phase === "final" && (
                    <span className="ml-auto rounded-full bg-tint-green/15 px-2.5 py-0.5 text-caption-2 font-bold text-tint-green">
                        FINAL
                    </span>
                )}
            </div>
            <p className="text-body text-label-primary">
                {locale === "zh" && session.nl_summary_zh ? session.nl_summary_zh : session.nl_summary}
            </p>
            {observationText && (
                <div className="mt-2 rounded-lg bg-surface-tertiary/60 p-3">
                    <p className="text-caption-1 text-label-tertiary italic">
                        💭 {observationText}
                    </p>
                </div>
            )}
            <div className="mt-2 flex gap-3 text-caption-2 text-label-quaternary">
                <span>{session.model_id}</span>
                <span>{session.input_tokens + session.output_tokens} tokens</span>
                <span>{session.latency_ms}ms</span>
            </div>
        </div>
    );
}

export default async function DailyPicksPage() {
    const locale = await getLocaleFromCookies();

    let todayPick: Awaited<ReturnType<typeof getDailyPick>> | null = null;
    try {
        todayPick = await getDailyPick();
    } catch {
        // no pick today
    }

    let history: Awaited<ReturnType<typeof getDailyPickHistory>> = [];
    try {
        history = await getDailyPickHistory(30);
    } catch {
        // no history
    }

    let sessions: NovaSession[] = [];
    try {
        sessions = await getNovaSessions();
    } catch {
        // no sessions
    }

    const hasFinal = sessions.some((s) => s.phase === "final");

    return (
        <section className="space-y-6">
            <h1 className="text-title-1 font-bold text-label-primary opacity-0 animate-slide-up stagger-1">
                {locale === "zh" ? "🏆 每日推荐交易者" : "🏆 Daily Pick"}
            </h1>

            {/* ── Nova Thinking Timeline ── */}
            {sessions.length > 0 && (
                <div className="opacity-0 animate-slide-up stagger-2">
                    <div className="flex items-center gap-2 mb-4">
                        <Brain className="h-5 w-5 text-tint-purple" />
                        <h2 className="text-headline font-bold text-label-primary">
                            {locale === "zh" ? "Nova 分析过程" : "Nova's Analysis Timeline"}
                        </h2>
                        {!hasFinal && (
                            <span className="ml-2 flex items-center gap-1 rounded-full bg-tint-blue/10 px-3 py-1 text-caption-2 font-medium text-tint-blue">
                                <Clock className="h-3 w-3" />
                                {locale === "zh" ? "分析中..." : "Analyzing..."}
                            </span>
                        )}
                    </div>

                    <div className="relative space-y-3 pl-6">
                        {/* Vertical timeline line */}
                        <div className="absolute left-2 top-2 bottom-2 w-px bg-separator" />

                        {sessions.map((session) => (
                            <div key={session.id} className="relative">
                                {/* Timeline dot */}
                                <div className={[
                                    "absolute -left-[18px] top-5 h-2.5 w-2.5 rounded-full",
                                    session.phase === "final" ? "bg-tint-green" : "bg-tint-blue",
                                ].join(" ")} />
                                <SessionCard session={session} locale={locale} />
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* ── Today's Pick (from final) ── */}
            {todayPick ? (
                <div className="rounded-2xl bg-gradient-to-br from-tint-blue/10 via-surface-secondary to-tint-purple/10 p-6 shadow-elevation-1 opacity-0 animate-slide-up stagger-3">
                    <div className="flex items-center gap-2 mb-4">
                        <Sparkles className="h-5 w-5 text-tint-blue" />
                        <span className="text-headline font-bold text-label-primary">
                            {locale === "zh" ? "今日推荐" : "Today's Recommendation"}
                        </span>
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                        <div className="rounded-xl bg-surface-primary/60 p-4">
                            <span className="text-caption-1 text-label-tertiary block mb-1">
                                {locale === "zh" ? "交易者" : "Trader"}
                            </span>
                            <Link
                                href={`/wallets/${todayPick.pick.wallet_id}`}
                                className="text-subheadline font-mono font-semibold text-tint-blue hover:underline"
                            >
                                {todayPick.wallet?.address
                                    ? `${todayPick.wallet.address.slice(0, 6)}…${todayPick.wallet.address.slice(-4)}`
                                    : `#${todayPick.pick.wallet_id}`}
                            </Link>
                        </div>

                        <div className="rounded-xl bg-surface-primary/60 p-4">
                            <span className="text-caption-1 text-label-tertiary block mb-1">Smart Score</span>
                            <span className="text-title-2 font-bold text-tint-green">{todayPick.pick.smart_score}</span>
                        </div>

                        <div className="rounded-xl bg-surface-primary/60 p-4">
                            <span className="text-caption-1 text-label-tertiary block mb-1">
                                {locale === "zh" ? "已实现盈亏" : "Realized PnL"}
                            </span>
                            <span className={`text-title-2 font-bold ${todayPick.pick.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
                                {todayPick.pick.realized_pnl >= 0 ? "+" : ""}{todayPick.pick.realized_pnl.toFixed(2)}
                            </span>
                        </div>

                        <div className="rounded-xl bg-surface-primary/60 p-4">
                            <span className="text-caption-1 text-label-tertiary block mb-1">
                                {locale === "zh" ? "交易次数" : "Trades"}
                            </span>
                            <span className="text-title-2 font-bold text-label-primary">{todayPick.pick.total_trades}</span>
                        </div>
                    </div>

                    {todayPick.pick.reason_summary && (
                        <div className="mt-4 rounded-xl bg-surface-primary/60 p-4">
                            <span className="text-caption-1 text-label-tertiary block mb-1">
                                {locale === "zh" ? "推荐理由" : "Recommendation Reason"}
                            </span>
                            <p className="text-body text-label-secondary">
                                {locale === "zh" && todayPick.pick.reason_summary_zh
                                    ? todayPick.pick.reason_summary_zh
                                    : todayPick.pick.reason_summary}
                            </p>
                        </div>
                    )}
                </div>
            ) : (
                <div className="rounded-2xl bg-surface-secondary p-6 text-center opacity-0 animate-slide-up stagger-3">
                    <p className="text-body text-label-tertiary">
                        {sessions.length > 0
                            ? (locale === "zh" ? "Nova 还在分析中，稍后会做出最终推荐..." : "Nova is still analyzing. Final pick coming soon...")
                            : (locale === "zh" ? "今日推荐尚未生成，请稍后再来。" : "Today's pick has not been generated yet. Check back later.")}
                    </p>
                </div>
            )}

            {/* ── History ── */}
            {history.length > 0 && (
                <div className="opacity-0 animate-slide-up stagger-4">
                    <h2 className="text-title-3 font-bold text-label-primary mb-4">
                        {locale === "zh" ? "推荐历史 & 跟单结果" : "History & Follow Results"}
                    </h2>

                    <div className="space-y-3">
                        {history.map((pick) => (
                            <div
                                key={pick.id}
                                className="flex items-center justify-between rounded-xl bg-surface-secondary px-5 py-4 shadow-elevation-1 transition-colors hover:bg-surface-tertiary/50"
                            >
                                <div className="flex items-center gap-4 min-w-0">
                                    <span className="text-caption-1 tabular-nums text-label-tertiary whitespace-nowrap">
                                        {pick.pick_date.slice(0, 10)}
                                    </span>
                                    <Link
                                        href={`/wallets/${pick.wallet_id}`}
                                        className="text-subheadline font-semibold text-tint-blue hover:underline truncate"
                                    >
                                        Wallet #{pick.wallet_id}
                                    </Link>
                                    <span className="text-caption-1 text-label-tertiary hidden sm:block">
                                        Score: {pick.smart_score}
                                    </span>
                                </div>

                                <div className="flex items-center gap-4 shrink-0">
                                    {pick.follow_pnl != null ? (
                                        <div className="flex items-center gap-1.5">
                                            {pick.follow_pnl >= 0 ? (
                                                <TrendingUp className="h-4 w-4 text-tint-green" />
                                            ) : (
                                                <TrendingDown className="h-4 w-4 text-tint-red" />
                                            )}
                                            <span className={`text-subheadline font-bold tabular-nums ${pick.follow_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
                                                {pick.follow_pnl >= 0 ? "+" : ""}{pick.follow_pnl.toFixed(2)}
                                            </span>
                                            <span className="text-caption-2 text-label-quaternary">
                                                ({pick.trades_followed} {locale === "zh" ? "笔" : "trades"})
                                            </span>
                                        </div>
                                    ) : (
                                        <span className="text-caption-1 text-label-quaternary">
                                            {locale === "zh" ? "待结算" : "Pending"}
                                        </span>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </section>
    );
}
