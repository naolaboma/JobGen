"use client";
import React from "react";
import { Menu, X } from "lucide-react";
import { Button } from "@/app/components/landingPage/ui/button";
import Link from "next/link";
import { useSession, signOut } from "next-auth/react";
import { apiUrl } from "@/lib/api";
import Image from "next/image";

export function LandingHeader({
  mobileMenuOpen,
  setMobileMenuOpen,
}: {
  mobileMenuOpen: boolean;
  setMobileMenuOpen: (open: boolean) => void;
}) {
  const { data: session, status } = useSession();
  const isAuth = status === "authenticated";

  async function handleLogout() {
    try {
      // Try to inform backend
      if ((session as any)?.accessToken) {
        await fetch(apiUrl("/api/v1/auth/logout"), {
          method: "POST",
          headers: { Authorization: `Bearer ${(session as any).accessToken}` },
        }).catch(() => {});
      }
    } finally {
      await signOut({ callbackUrl: "/" });
    }
  }

  return (
    <header className="sticky top-0 z-50 bg-white text-gray-900 border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-3 group">
            <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-sm ring-1 ring-gray-200 group-hover:scale-[1.03] transition-transform">
              <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center">
                <div className="flex gap-1">
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                </div>
              </div>
            </div>
            <span className="text-xl font-bold">JobGen</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link href="/jobs" className="hover:text-[#44C3BB] transition-colors">
              Jobs
            </Link>
            {isAuth && (
              <>
                <Link
                  href="/user-home/personalized-jobs"
                  className="hover:text-[#44C3BB] transition-colors"
                >
                  My Matches
                </Link>
                <Link
                  href="/cv"
                  className="hover:text-[#44C3BB] transition-colors"
                >
                  Upload CV
                </Link>
              </>
            )}
            <Link
              href="/contact"
              className="hover:text-[#44C3BB] transition-colors"
            >
              Contact
            </Link>
          </nav>

          {/* Desktop Right Side */}
          <div className="hidden md:flex items-center space-x-4">
            {isAuth ? (
              <div className="flex items-center gap-3">
                <Link href="/user-home" className="hover:text-[#44C3BB]">
                  Dashboard
                </Link>
                <Link
                  href="/profile"
                  className="rounded-full bg-white w-8 h-8 flex items-center justify-center overflow-hidden ring-2 ring-gray-200 hover:ring-[#44C3BB]/50"
                >
                  <Image
                    src="/professional-woman-dark-hair.png"
                    alt="Profile"
                    width={32}
                    height={32}
                    className="object-cover"
                  />
                </Link>
                <div className="px-3 py-1.5 rounded-full bg-gray-100 text-gray-800 text-sm">
                  {(session?.user?.email || "").split("@")[0] || "user"}
                </div>
                <Button
                  onClick={handleLogout}
                  className="bg-[#44C3BB] text-white hover:bg-[#3bb3ac]"
                >
                  Sign Out
                </Button>
              </div>
            ) : (
              <>
                <Link href="/login" passHref>
                  <Button
                    variant="outline"
                    className="text-gray-900 border-gray-300 hover:bg-gray-50"
                  >
                    Sign In
                  </Button>
                </Link>
                <Link href="/register" passHref>
                  <Button className="bg-[#44C3BB] hover:bg-[#1e7a73] text-white">
                    Get Started
                  </Button>
                </Link>
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden p-2 hover:bg-gray-100 rounded-md"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200 bg-white">
            <div className="flex flex-col space-y-4">
              <Link href="/jobs" className="hover:text-[#44C3BB]">
                Jobs
              </Link>
              {isAuth && (
                <>
                  <Link
                    href="/user-home/personalized-jobs"
                    className="hover:text-[#44C3BB]"
                  >
                    My Matches
                  </Link>
                  <Link href="/cv" className="hover:text-[#44C3BB]">
                    Upload CV
                  </Link>
                </>
              )}
              <Link href="/contact" className="hover:text-[#44C3BB]">
                Contact
              </Link>
              <div className="flex flex-col space-y-2 pt-4 border-t border-gray-200">
                {isAuth ? (
                  <>
                    <Link href="/user-home" className="hover:text-[#44C3BB]">
                      Dashboard
                    </Link>
                    <Button
                      onClick={handleLogout}
                      className="w-full bg-[#44C3BB] hover:bg-[#3AB5AD] text-white"
                    >
                      Sign Out
                    </Button>
                  </>
                ) : (
                  <>
                    <Link href="/login" passHref>
                      <Button
                        variant="outline"
                        className="w-full border-gray-300 hover:bg-gray-50"
                      >
                        Sign In
                      </Button>
                    </Link>
                    <Link href="/register" passHref>
                      <Button className="w-full bg-[#44C3BB] hover:bg-[#3AB5AD] text-white">
                        Get Started
                      </Button>
                    </Link>
                  </>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
