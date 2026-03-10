"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface WalletInputProps {
  onSubmit: (input: string) => void;
  isLoading: boolean;
}

export function WalletInput({ onSubmit, isLoading }: WalletInputProps) {
  const [value, setValue] = useState("");
  const [error, setError] = useState("");

  const ETH_ADDRESS = /^0x[a-fA-F0-9]{40}$/;
  const PROFILE_URL = /polymarket\.com\/profile\/(0x[a-fA-F0-9]{40})/;

  const validate = (input: string): boolean => {
    const trimmed = input.trim();
    if (!trimmed) {
      setError("Please enter a wallet address or Polymarket profile URL");
      return false;
    }
    if (ETH_ADDRESS.test(trimmed) || PROFILE_URL.test(trimmed)) {
      setError("");
      return true;
    }
    setError(
      "Must be an Ethereum address (0x followed by 40 hex chars) or a Polymarket profile URL (polymarket.com/profile/0x...)"
    );
    return false;
  };

  const handleSubmit = () => {
    if (validate(value)) {
      onSubmit(value.trim());
    }
  };

  return (
    <div className="w-full max-w-xl space-y-4">
      <div className="flex gap-3">
        <Input
          placeholder="0x... or polymarket.com/profile/0x..."
          value={value}
          onChange={(e) => {
            setValue(e.target.value);
            if (error) setError("");
          }}
          onKeyDown={(e) => e.key === "Enter" && handleSubmit()}
          disabled={isLoading}
          className="h-12 text-base bg-white/10 border-white/20 text-white placeholder:text-white/40"
        />
        <Button
          onClick={handleSubmit}
          disabled={isLoading || !value.trim()}
          className="h-12 px-8 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white font-semibold"
        >
          {isLoading ? (
            <span className="flex items-center gap-2">
              <span className="h-4 w-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              Analyzing
            </span>
          ) : (
            "Analyze"
          )}
        </Button>
      </div>
      {error && <p className="text-red-400 text-sm">{error}</p>}
    </div>
  );
}
