export const runtime = "nodejs";

import { NextRequest, NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/authOptions";

function backendBase() {
  const raw = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  return raw.replace(/\/+$/, "");
}

export async function POST(req: NextRequest) {
  try {
    const form = await req.formData();
    const file = form.get("file") as unknown as File | null;

    if (!file) {
      return NextResponse.json({ error: "file is required" }, { status: 400 });
    }

    const session = await getServerSession(authOptions as any);
    const token = (session as any)?.accessToken as string | undefined;
    const incomingAuth = req.headers.get("authorization") || undefined;
    const authHeader = token ? `Bearer ${token}` : incomingAuth;

    const backend = backendBase();

    // If BACKEND_URL/NEXT_PUBLIC_API_URL is present, attempt to proxy to backend
    if (backend) {
      try {
        const f = new FormData();
        f.append("file", file as unknown as Blob, (file as any).name || "cv.pdf");

        const r = await fetch(`${backend}/api/v1/cv/parse`, {
          method: "POST",
          headers: {
            ...(authHeader ? { Authorization: authHeader } : {}),
            Accept: "application/json",
          },
          body: f,
        });

        const j = await r.json().catch(() => ({}));
        return NextResponse.json(j, { status: r.status });
      } catch (e: any) {
        // Fall through to mock below
      }
    }

    // Mock response for local/dev without backend
    const mockJobId = `mock-${Math.random().toString(36).slice(2, 8)}`;
    return NextResponse.json(
      {
        message: "Mock: CV parsing job accepted.",
        jobId: mockJobId,
      },
      { status: 202 }
    );
  } catch (err: any) {
    return NextResponse.json({ error: err?.message || "Unexpected error" }, { status: 500 });
  }
}
