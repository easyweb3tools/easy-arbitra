export function fallbackSummary(params: {
  strategyType?: string;
  smartScore?: number;
  tradeCount?: number;
  realizedPnl?: number;
  locale?: "en" | "zh";
}) {
  const {
    strategyType = "unknown",
    smartScore = 0,
    tradeCount = 0,
    realizedPnl = 0,
    locale = "en"
  } = params;

  let riskStatementEn = "showing elevated volatility, caution is advised";
  let riskStatementZh = "近期波动偏大，建议谨慎关注";
  if (smartScore >= 80) {
    riskStatementEn = "overall behavior is stable";
    riskStatementZh = "整体表现较稳定";
  } else if (smartScore >= 60) {
    riskStatementEn = "performance is acceptable and worth monitoring";
    riskStatementZh = "表现尚可，建议持续观察";
  }

  if (locale === "zh") {
    return `该钱包采用${strategyType}策略，当前评分${smartScore}，交易${tradeCount}次，已实现收益${realizedPnl.toFixed(2)}，${riskStatementZh}。`;
  }
  return `This wallet follows ${strategyType} style with score ${smartScore}, ${tradeCount} trades and realized PnL ${realizedPnl.toFixed(2)}; ${riskStatementEn}.`;
}
