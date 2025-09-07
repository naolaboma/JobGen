import React, { JSX } from "react";
import { Button } from "./ui/button";

export const FeaturesSection = (): JSX.Element => {
  return (
    <section className="flex w-full items-center justify-center gap-8 md:gap-12 lg:gap-[98px] px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px]">
      <div className="inline-flex flex-col lg:flex-row items-center justify-center gap-8 md:gap-12 lg:gap-[98px] relative flex-[0_0_auto] w-full max-w-[1480px]">
        <div className="inline-flex flex-col items-start gap-1 relative flex-[0_0_auto]">
          <div className="relative w-full max-w-[500px] lg:max-w-[714px] h-[300px] md:h-[400px] lg:h-[532.09px] bg-primary-100 rounded-lg" />
        </div>

        <div className="flex flex-col w-full lg:w-[669px] items-start gap-8 md:gap-12 lg:gap-[60px] relative">
          <div className="flex flex-col items-start gap-4 md:gap-6 relative self-stretch w-full flex-[0_0_auto]">
            <img
              className="absolute w-[400px] md:w-[500px] lg:w-[648px] h-[43px] md:h-[56px] lg:h-[70px] top-16 md:top-20 lg:top-28 left-1/2 lg:left-[-11px] transform -translate-x-1/2 lg:translate-x-0 hidden md:block"
              alt="Element"
              src="/element-9.png"
            />

            <h2 className="relative self-stretch mt-[-1.00px] font-header-h1 font-[number:var(--header-h1-font-weight)] text-[#212529] text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center lg:text-left tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
              Apply with Certainty
            </h2>

            <p className="relative self-stretch [font-family:'Inter',Helvetica] font-normal text-[#212529] text-base md:text-lg leading-6 md:leading-[18px] text-center lg:text-left">
              <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                Stop second-guessing your resume. Get the confidence that comes
                from knowing your CV is optimized to pass automated screens and
                impress human recruiters.
              </span>
            </p>
          </div>

          <Button className="gap-2.5 px-6 py-4 md:px-10 md:py-5 bg-[#7bbfb3] rounded-lg overflow-hidden inline-flex items-center justify-center relative flex-[0_0_auto] h-auto hover:bg-[#6ba89d] mx-auto lg:mx-0">
            <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-lg leading-[18px]">
              <span className="font-bold tracking-[-0.06px] leading-[23px]">
                Let&apos;s Go
              </span>
            </span>

            <img
              className="relative flex-[0_0_auto] mr-[-0.50px]"
              alt="Icon"
              src="/icon-17.svg"
            />
          </Button>
        </div>
      </div>
    </section>
  );
};
