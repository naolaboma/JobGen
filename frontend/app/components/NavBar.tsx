"use client";

import { useState } from "react";
import Image from "next/image";
import { Menu, X } from "lucide-react";

export default function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => setIsMobileMenuOpen((prev) => !prev);
  const closeMobileMenu = () => setIsMobileMenuOpen(false);

  return (
    <header className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center">
              <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center">
                <div className="flex gap-1">
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                  <div className="w-1 h-1 bg-white rounded-full"></div>
                </div>
              </div>
            </div>
            <span className="text-xl font-semibold text-white dark:text-gray-100">
              JobGen
            </span>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            <a
              href="#"
              className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
            >
              Jobs
            </a>
            <a
              href="#"
              className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
            >
              Blog
            </a>
            <a
              href="#"
              className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
            >
              About
            </a>
            <a
              href="#"
              className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
            >
              Contact
            </a>
            <div className="w-8 h-8 rounded-full overflow-hidden cursor-pointer ring-0 hover:ring-2 hover:ring-[#44C3BB]/60 hover:scale-150 transition-all duration-200 z-50">
              <Image
                src="/professional-woman-dark-hair.png"
                alt="Profile"
                className="w-full h-full object-cover"
                width={48}
                height={48}
              />
            </div>
          </nav>

          {/* Mobile Menu Button */}
          <button
            onClick={toggleMobileMenu}
            className="md:hidden p-2 rounded-md hover:bg-white/10 transition-colors cursor-pointer"
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
          <div className="md:hidden py-4 border-t border-white/20 dark:border-gray-700 bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800">
            <div className="flex flex-col gap-4">
              <a
                href="#"
                onClick={closeMobileMenu}
                className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
              >
                Blog
              </a>
              <a
                href="#"
                onClick={closeMobileMenu}
                className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
              >
                About
              </a>
              <a
                href="/contact-us"
                onClick={closeMobileMenu}
                className="hover:text-gray-200 dark:hover:text-gray-300 transition-colors"
              >
                Contact
              </a>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
