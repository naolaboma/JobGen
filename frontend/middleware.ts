import { withAuth } from "next-auth/middleware";

export default withAuth(
  // No custom middleware logic; rely on authorized() to gate access
  function middleware() {},
  {
    pages: {
      signIn: "/login",
    },
    callbacks: {
      // Only allow if there is a valid session token
      authorized: ({ token }) => !!token,
    },
  }
);

export const config = {
  matcher: [
    // Chat
    "/chat",
    "/chat/:path*",
    // User home (jobs, etc.)
    "/user-home",
    "/user-home/:path*",
    // Profile
    "/profile",
    "/profile/:path*",
    // Settings
    "/settings",
    "/settings/:path*",
    // Dashboard
    "/dashboard",
    "/dashboard/:path*",
    // Notifications (if needed)
    "/notifications",
    "/notifications/:path*",
  ],
};
