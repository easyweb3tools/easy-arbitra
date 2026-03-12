import { NextResponse } from "next/server";
import { analyzeWalletStream } from "@/lib/ai";

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

    const encoder = new TextEncoder();
    const stream = new ReadableStream({
      async start(controller) {
        try {
          await analyzeWalletStream(walletInput, (event) => {
            const line = `data: ${JSON.stringify(event)}\n\n`;
            controller.enqueue(encoder.encode(line));
          });
        } catch (error) {
          const errorEvent = {
            type: "error",
            data: error instanceof Error ? error.message : "Internal server error",
          };
          controller.enqueue(
            encoder.encode(`data: ${JSON.stringify(errorEvent)}\n\n`)
          );
        } finally {
          controller.close();
        }
      },
    });

    return new Response(stream, {
      headers: {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
      },
    });
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
