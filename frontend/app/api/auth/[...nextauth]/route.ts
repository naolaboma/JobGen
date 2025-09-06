import NextAuth from "next-auth";
import { authOptions } from "../../../../lib/authOptions"; // Correct import path

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
