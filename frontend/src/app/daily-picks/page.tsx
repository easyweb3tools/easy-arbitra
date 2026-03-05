import Link from "next/link";
import { Sparkles, TrendingUp, TrendingDown, Brain, Clock, Users } from "lucide-react";
import { getDailyPick, getDailyPickHistory, getNovaTimeline, getNovaDecisionExplanation, getNovaCandidates } from "@/lib/api";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { ThinkingTimeline } from "@/components/nova/ThinkingTimeline";
import { DecisionExplainer } from "@/components/nova/DecisionExplainer";
import { CandidateScoreCard } from "@/components/nova/CandidateScoreCard";

export const dynamic = "force-dynamic";

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

    // Get today's date for timeline
    const today = new Date().toISOString().split('T')[0];
    let timeline: Awaited<ReturnType<typeof getNovaTimeline>> = [];
    try {
        timeline = await getNovaTimeline(today);
    } catch {
        // no timeline
    }

    // Get candidates
    let candidates: Awaited<ReturnType<typeof getNovaCandidates>> = [];
    try {
        candidates = await getNovaCandidates(today);
    } catch {
        // no candidates
    }

    // Get decision explanation if we have a pick
    let explanation: Awaited<ReturnType<typeof getNovaDecisionExplanation>> | null = null;
    if (todayPick) {
        try {
            explanation = await getNovaDecisionExplanation(todayPick.pick.id);
        } catch {
            // no explanation
        }
    }

    const hasFinal = timeline.some((r) => r.session.phase === "final");

    return (
        <section className="space-y-6">
            <h1 className="text-title-1 font-bold text-label-primary opacity-0 animate-slide-up stagger-1">
                {locale === "zh" ? "🏆 每日推荐交易者" : "🏆 Daily Pick"}
            </h1>

            {/* ── Nova Thinking Timeline ── */}
            {timeline.length > 0 && (
                <div className="opacity-0 animate-slide-up stagger-2">
                    <div className="flex items-center gap-2 mb-4">
                        <Brain className="h-5 w-5 text-tint-purple" />
                        <h2 className="text-headline font-bold text-label-primary">
                            {locale === "zh" ? "Nova 思考过程" : "Nova's Thinking Process"}
                        </h2>
                        {!hasFinal && (
                            <span className="ml-2 flex items-center gap-1 rounded-full bg-tint-blue/10 px-3 py-1 text-caption-2 font-medium text-tint-blue">
                                <Clock className="h-3 w-3" />
                                {locale === "zh" ? "分析中..." : "Analyzing..."}
                            </span>
                        )}
                    </div>

                    <ThinkingTimeline rounds={timeline} locale={locale} />
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

                    {/* Decision Explanation */}
                    {explanation && (
                        <div className="mt-4 rounded-xl bg-surface-primary/60 p-4">
                            <h3 className="text-subheadline font-bold text-label-primary mb-4">
                                {locale === "zh" ? "Nova 决策分析" : "Nova's Decision Analysis"}
                            </h3>
                            <DecisionExplainer explanation={explanation} locale={locale} />
                        </div>
                    )}
                </div>
            ) : (
                <div className="rounded-2xl bg-surface-secondary p-6 text-center opacity-0 animate-slide-up stagger-3">
                    <p className="text-body text-label-tertiary">
                        {timeline.length > 0
                            ? (locale === "zh" ? "Nova 还在分析中，稍后会做出最终推荐..." : "Nova is still analyzing. Final pick coming soon...")
                            : (locale === "zh" ? "今日推荐尚未生成，请稍后再来。" : "Today's pick has not been generated yet. Check back later.")}
                    </p>
                </div>
            )}

            {/* ── Candidate Wallets ── */}
            {candidates.length > 0 && (
                <div className="opacity-0 animate-slide-up stagger-4">
                    <div className="flex items-center gap-2 mb-4">
                        <Users className="h-5 w-5 text-tint-blue" />
                        <h2 className="text-headline font-bold text-label-primary">
                            {locale === "zh" ? "候选钱包评分" : "Candidate Wallet Scores"}
                        </h2>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {candidates.slice(0, 6).map((candidate) => (
                            <CandidateScoreCard
                                key={candidate.wallet_id}
                                candidate={candidate}
                                locale={locale}
                            />
                        ))}
                    </div>

                    {candidates.length > 6 && (
                        <div className="mt-4 text-center">
                            <p className="text-caption-1 text-label-tertiary">
                                {locale === "zh" 
                                    ? `还有 ${candidates.length - 6} 个候选钱包...` 
                                    : `${candidates.length - 6} more candidates...`}
                            </p>
                        </div>
                    )}
                </div>
            )}

            {/* ── History ── */}
            {history.length > 0 && (
                <div className="opacity-0 animate-slide-up stagger-5">
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
