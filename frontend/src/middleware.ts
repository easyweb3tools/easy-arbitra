// No-op middleware — all routes are public in the simplified product.
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(_request: NextRequest) {
    return NextResponse.next();
}
