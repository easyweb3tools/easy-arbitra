"use client";

import {
  Radar,
  RadarChart as RechartsRadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  ResponsiveContainer,
} from "recharts";
import type { RadarChartData } from "@/lib/types";

interface RadarChartProps {
  data: RadarChartData;
}

export function RadarChart({ data }: RadarChartProps) {
  const chartData = [
    { axis: "Entry Timing", value: data.entry_timing, fullMark: 1 },
    { axis: "Position Size", value: data.size_ratio, fullMark: 1 },
    { axis: "ROI", value: data.roi, fullMark: 1 },
  ];

  return (
    <div className="w-full h-[280px]">
      <ResponsiveContainer width="100%" height="100%">
        <RechartsRadarChart data={chartData} cx="50%" cy="50%" outerRadius="70%">
          <PolarGrid stroke="rgba(255,255,255,0.1)" />
          <PolarAngleAxis
            dataKey="axis"
            tick={{ fill: "rgba(255,255,255,0.7)", fontSize: 12 }}
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 1]}
            tick={{ fill: "rgba(255,255,255,0.4)", fontSize: 10 }}
            tickCount={5}
          />
          <Radar
            name="Style"
            dataKey="value"
            stroke="#8b5cf6"
            fill="#8b5cf6"
            fillOpacity={0.3}
            strokeWidth={2}
          />
        </RechartsRadarChart>
      </ResponsiveContainer>
    </div>
  );
}
