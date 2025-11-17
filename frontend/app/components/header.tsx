"use client";

import { useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { Menu, X } from "lucide-react";
import { useSession } from "next-auth/react";

export default function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const { status } = useSession();
  const isAuth = status === "authenticated";

  const toggleMobileMenu = () => setIsMobileMenuOpen((prev) => !prev);
  const closeMobileMenu = () => setIsMobileMenuOpen(false);

  return (
    <header className="bg-white text-gray-900 border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Brand */}
          <Link href="/" className="flex items-center gap-3 group">
            <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-sm ring-1 ring-gray-200 group-hover:scale-[1.03] transition-transform">
              <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center" />
            </div>
            <span className="text-xl font-semibold">JobGen</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            <Link
              href="/jobs"
              className="hover:text-[#44C3BB] transition-colors"
            >
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
            <Link
              href="/profile"
              className="w-8 h-8 rounded-full overflow-hidden ring-2 ring-gray-200 hover:ring-[#44C3BB]/50 transition-shadow"
            >
              <Image
                src="/professional-woman-dark-hair.png"
                alt="Profile"
                className="w-full h-full object-cover"
                width={48}
                height={48}
                priority
              />
            </Link>
          </nav>

          {/* Mobile Menu Button */}
          <button
            onClick={toggleMobileMenu}
            className="md:hidden p-2 rounded-md hover:bg-gray-100 transition-colors"
            aria-label="Toggle Menu"
          >
            {isMobileMenuOpen ? (
              <X className="w-6 h-6" />
            ) : (
              <Menu className="w-6 h-6" />
            )}
          </button>
        </div>

        {/* Mobile Navigation */}
        {isMobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200">
            <div className="flex flex-col gap-4">
              <Link
                href="/jobs"
                onClick={closeMobileMenu}
                className="hover:text-[#44C3BB] transition-colors"
              >
                Jobs
              </Link>
              {isAuth && (
                <>
                  <Link
                    href="/user-home/personalized-jobs"
                    onClick={closeMobileMenu}
                    className="hover:text-[#44C3BB] transition-colors"
                  >
                    My Matches
                  </Link>
                  <Link
                    href="/cv"
                    onClick={closeMobileMenu}
                    className="hover:text-[#44C3BB] transition-colors"
                  >
                    Upload CV
                  </Link>
                </>
              )}
              <Link
                href="/contact"
                onClick={closeMobileMenu}
                className="hover:text-[#44C3BB] transition-colors"
              >
                Contact
              </Link>
              <Link
                href="/profile"
                onClick={closeMobileMenu}
                className="flex items-center gap-3 hover:text-[#44C3BB] transition-colors"
              >
                <div className="w-8 h-8 rounded-full overflow-hidden ring-2 ring-gray-200">
                  <Image
                    src="/professional-woman-dark-hair.png"
                    alt="Profile"
                    width={32}
                    height={32}
                    className="object-cover"
                  />
                </div>
                <span>Profile</span>
              </Link>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
