"use client";

import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { WalletCard as WalletCardType } from "@/lib/types";

interface WalletCardProps {
  data: WalletCardType;
}

export function WalletCard({ data }: WalletCardProps) {
  const shortAddress = data.address
    ? `${data.address.slice(0, 6)}...${data.address.slice(-4)}`
    : "—";

  return (
    <Card className="bg-white/5 border-white/10">
      <CardHeader className="pb-3">
        <div className="flex items-center gap-3">
          {data.profile_image ? (
            <img
              src={data.profile_image}
              alt={data.display_name}
              className="w-12 h-12 rounded-full border border-white/20"
            />
          ) : (
            <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg">
              {(data.display_name || "?")[0].toUpperCase()}
            </div>
          )}
          <div>
            <h3 className="font-semibold text-white text-lg">
              {data.display_name || "Anonymous"}
            </h3>
            <p className="text-xs text-white/50 font-mono">{shortAddress}</p>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-2">
          <Badge className="bg-orange-500/20 text-orange-300 border-orange-500/30">
            {data.sport}
          </Badge>
          <span className="text-sm text-white/60">
            {data.total_trades} trades analyzed
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
