import { ArrowRightIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";

export const CareerOpportunitiesSection = (): JSX.Element => {
  return (
    <section className="bg-[#7bbfb3] py-16 md:py-24 lg:py-[140px] px-4 md:px-8 lg:px-[220px] w-full relative">
      <div className="flex flex-col items-center justify-center gap-8 md:gap-12 lg:gap-[60px] max-w-[1481px] mx-auto">
        <div className="flex flex-col items-center gap-4 md:gap-6 max-w-[1064px] w-full">
          <div className="text-center">
            <h2 className="font-header-h1 font-[number:var(--header-h1-font-weight)] text-white text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)] mb-4">
              Ready to Unlock Global Opportunities?
            </h2>

            <img
              className="w-[200px] md:w-[280px] lg:w-[314px] h-6 md:h-8 lg:h-9 mx-auto mb-4"
              alt="Element"
              src="/element-5.png"
            />
          </div>

          <p className="[font-family:'Inter',Helvetica] font-bold text-white text-base md:text-lg text-center tracking-[-0.06px] leading-6 md:leading-[30px] max-w-[1064px] px-4">
            Stop searching endlessly. Start matching intelligently. Join the
            waitlist for early access and be the first to launch your
            international career.
          </p>
        </div>

        <Button className="bg-[#4f9cf9] hover:bg-[#4f9cf9]/90 rounded-lg px-6 py-4 md:px-10 md:py-5 h-auto relative overflow-hidden">
          <img
            className="absolute w-[300px] md:w-[400px] lg:w-[475px] h-[500px] md:h-[700px] lg:h-[836px] top-[-350px] md:top-[-450px] lg:top-[-555px] left-[-500px] md:left-[-700px] lg:left-[-837px] pointer-events-none"
            alt="Background"
            src="/background-1.png"
          />

          <span className="[font-family:'Inter',Helvetica] font-bold text-white text-base md:text-lg tracking-[-0.06px] leading-[23px] relative z-10">
            Get Early Access
          </span>

          <ArrowRightIcon className="w-3 h-3 md:w-4 md:h-4 ml-1 md:ml-2 relative z-10" />
        </Button>
      </div>
    </section>
  );
};
