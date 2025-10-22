import { withAuth } from "next-auth/middleware";

export default withAuth({
  pages: {
    signIn: "/login",
  },
});

export const config = {
  matcher: [
    "/chat/:path*",
    "/user-home/:path*",
    "/profile/:path*",
    "/settings/:path*",
    "/dashboard/:path*",
  ],
};
