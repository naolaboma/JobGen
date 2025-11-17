export const runtime = "nodejs";

import { NextRequest, NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/authOptions";

function backendBase() {
  const raw =
    process.env.BACKEND_URL ||
    process.env.NEXT_PUBLIC_API_URL ||
    "http://localhost:8080";
  return raw.replace(/\/+$/, "");
}

export async function GET(
  req: NextRequest,
  context: { params: Promise<{ id: string }> }
) {
  const { id } = await context.params;
  try {
    const backend = backendBase();

    // Prefer server-session token; fallback to inbound Authorization header (for programmatic calls)
    const session = await getServerSession(authOptions as any);
    const token = (session as any)?.accessToken as string | undefined;
    const incomingAuth = req.headers.get("authorization") || undefined;
    const authHeader = token ? `Bearer ${token}` : incomingAuth;

    if (backend) {
      try {
        const r = await fetch(
          `${backend}/api/v1/cv/${encodeURIComponent(id)}`,
          {
            method: "GET",
            headers: {
              ...(authHeader ? { Authorization: authHeader } : {}),
              Accept: "application/json",
            },
          }
        );
        const j = await r.json().catch(() => ({}));
        return NextResponse.json(j, { status: r.status });
      } catch (_) {
        // Fall through to mock
      }
    }

    // Mock completed CV response
    const mock = {
      id,
      userId: "user-mock",
      fileStorageId: "fs-mock",
      fileName: "cv.pdf",
      status: "Completed",
      processingError: "",
      rawText: "",
      profileSummary: "Seasoned developer with 5+ years in TypeScript and Go.",
      experiences: [],
      educations: [],
      skills: ["TypeScript", "React", "Node.js", "Go"],
      suggestions: [
        {
          id: "s1",
          type: "quantification",
          content: "Add measurable outcomes to your recent role.",
          applied: false,
        },
        {
          id: "s2",
          type: "missing_keywords",
          content: "Consider adding Docker and CI/CD if applicable.",
          applied: false,
        },
      ],
      score: 78,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    return NextResponse.json(mock, { status: 200 });
  } catch (err: any) {
    return NextResponse.json(
      { error: err?.message || "Unexpected error" },
      { status: 500 }
    );
  }
}
