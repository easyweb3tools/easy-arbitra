"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import type { ReportData } from "@/lib/types";

interface ReportSummaryProps {
  report: ReportData;
  explanation: string;
}

export function ReportSummary({ report, explanation }: ReportSummaryProps) {
  return (
    <Card className="bg-white/5 border-white/10">
      <CardHeader className="pb-3">
        <CardTitle className="text-white/80 text-sm font-medium">
          Style Analysis
        </CardTitle>
        <div className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
          {report.style_label}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-white/60">{report.summary_context}</p>
        <Separator className="bg-white/10" />
        <div className="bg-white/5 rounded-lg p-4 border-l-2 border-purple-500">
          <p className="text-xs text-purple-300 mb-2 font-medium">
            Nova AI Analysis
          </p>
          <p className="text-sm text-white/80 whitespace-pre-wrap leading-relaxed">
            {explanation}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
