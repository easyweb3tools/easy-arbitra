import { NextResponse } from "next/server";
import { analyzeWallet } from "@/lib/bedrock";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { walletInput } = body;

    if (!walletInput || typeof walletInput !== "string") {
      return NextResponse.json(
        { error: "walletInput is required" },
        { status: 400 }
      );
    }

    const result = await analyzeWallet(walletInput);
    return NextResponse.json(result);
  } catch (error) {
    console.error("Analyze error:", error);
    return NextResponse.json(
      {
        error:
          error instanceof Error ? error.message : "Internal server error",
      },
      { status: 500 }
    );
  }
}
