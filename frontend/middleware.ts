import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { getToken } from "next-auth/jwt";

const PROTECTED_PREFIXES = [
  "/chat",
  "/user-home",
  "/profile",
  "/settings",
  "/dashboard",
  "/notifications",
];

export async function middleware(req: NextRequest) {
  const { pathname, search } = req.nextUrl;

  // Skip internal/static and auth API
  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/static") ||
    pathname.startsWith("/api/auth")
  ) {
    return NextResponse.next();
  }

  const token = await getToken({ req, secret: process.env.NEXTAUTH_SECRET });
  // Safely normalize exp to a number (ms)
  const expSec =
    typeof token?.exp === "number"
      ? token.exp
      : typeof token?.exp === "string"
      ? parseInt(token.exp, 10)
      : undefined;
  const expMs =
    typeof expSec === "number" && Number.isFinite(expSec)
      ? expSec * 1000
      : undefined;
  const isExpired = typeof expMs === "number" ? Date.now() >= expMs : false;
  const isAuth = !!token && !isExpired;
  const isProtected = PROTECTED_PREFIXES.some(
    (p) => pathname === p || pathname.startsWith(`${p}/`)
  );
  const isAuthPage = pathname === "/login" || pathname === "/register";

  // Redirect logged-in users away from auth pages
  if (isAuth && isAuthPage) {
    const url = req.nextUrl.clone();
    url.pathname = "/user-home/fallback-page";
    url.search = "";
    return NextResponse.redirect(url);
  }

  // Redirect unauthenticated or expired sessions to login with callbackUrl
  if (!isAuth && isProtected) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    const callbackUrl = encodeURIComponent(
      `${req.nextUrl.origin}${pathname}${search}`
    );
    url.search = `?callbackUrl=${callbackUrl}`;
    return NextResponse.redirect(url);
  }

  const res = NextResponse.next();
  if (isProtected) {
    res.headers.set(
      "Cache-Control",
      "private, no-store, max-age=0, must-revalidate"
    );
    res.headers.set("Vary", "Cookie");
  }
  return res;
}

export const config = {
  matcher: [
    "/chat",
    "/chat/:path*",
    "/user-home",
    "/user-home/:path*",
    "/profile",
    "/profile/:path*",
    "/settings",
    "/settings/:path*",
    "/dashboard",
    "/dashboard/:path*",
    "/notifications",
    "/notifications/:path*",
    "/login",
    "/register",
  ],
};
