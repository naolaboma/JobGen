import Link from "next/link";

export default function Footer() {
  return (
    <footer className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-200 w-full mt-16 flex-shrink-0">
      <div className="max-w-7xl mx-auto pb-8 pt-8 px-4 sm:px-7 lg:px-10">
        <div className="grid grid-cols-1 md:grid-cols-3 items-center gap-6 md:gap-0 text-center md:text-left">
          <div className="flex justify-center md:justify-start gap-8 mb-4 md:mb-0">
            <Link
              href="/privacy"
              className="hover:text-gray-200 dark:hover:text-white transition-colors"
            >
              Privacy
            </Link>
            <Link
              href="/terms"
              className="hover:text-gray-200 dark:hover:text-white transition-colors"
            >
              Terms
            </Link>
          </div>
          <div className="flex justify-center mb-4 md:mb-0">
            <p className="text-white dark:text-gray-200">
              © {new Date().getFullYear()} JobGen
            </p>
          </div>
          <div className="flex justify-center md:justify-end">
            <Link
              href="/contact"
              className="hover:text-gray-200 dark:hover:text-white transition-colors"
            >
              Contact
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
