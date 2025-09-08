import React from "react";
import { Button } from "@/app/components/landingPage/ui/button";
import Link from "next/link";

export function CtaSection() {
  return (
    <section className="py-20 bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800">
      <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
        <h2 className="text-3xl sm:text-4xl font-bold text-white dark:text-gray-100 mb-6">
          Ready to accelerate your career?
        </h2>
        <p className="text-xl text-teal-100 dark:text-gray-300 mb-8">
          Join thousands of professionals who have transformed their careers
          with JobGen.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link href="/register" passHref>
            <Button size="lg">Start Free Trial</Button>
          </Link>
          <Button size="lg" variant="outline">
            Schedule Demo
          </Button>
        </div>
      </div>
    </section>
  );
}
