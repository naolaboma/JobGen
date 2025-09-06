import { ChevronRightIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";

export const HeroSection = (): JSX.Element => {
  return (
    <section className="flex w-full items-center justify-center px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px] relative bg-[#7bbfb3]">
      <img
        className="absolute w-full h-[300px] md:h-[400px] lg:h-[547px] top-24 md:top-32 lg:top-48 left-0"
        alt="Element"
        src="/element.png"
      />

      <div className="flex flex-col lg:flex-row items-center justify-center relative w-full max-w-[1480px] gap-8 lg:gap-0">
        <div className="flex flex-col w-full lg:w-[656px] items-start gap-8 md:gap-12 lg:gap-[60px] relative">
          <div className="flex flex-col items-start gap-4 md:gap-6 relative self-stretch w-full">
            <h1 className="relative self-stretch mt-[-1.00px] font-header-h2 font-[number:var(--header-h2-font-weight)] text-white text-3xl md:text-5xl lg:text-[length:var(--header-h2-font-size)] tracking-[var(--header-h2-letter-spacing)] leading-tight md:leading-[var(--header-h2-line-height)] [font-style:var(--header-h2-font-style)]">
              Land Your Dream Remote Tech Job
            </h1>

            <p className="relative self-stretch [font-family:'Inter',Helvetica] font-normal text-white text-base md:text-lg leading-6 md:leading-[18px]">
              <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                JobGen is your AI career coach. Get personalized CV feedback and
                discover remote opportunities that match your skills—all in one
                place. Built for African tech talent.
              </span>
            </p>
          </div>

          <Button className="inline-flex items-center gap-2.5 px-6 py-4 md:p-5 h-auto bg-[#38b2ac] hover:bg-[#319795] rounded-lg">
            <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-lg leading-[18px]">
              <span className="font-bold tracking-[-0.06px] leading-[23px]">
                See How It Works
              </span>
            </span>

            <ChevronRightIcon className="relative w-[11px] h-[11px] mr-[-0.50px] text-white" />
          </Button>
        </div>

        <div className="flex flex-col items-start gap-1 relative flex-1 w-full lg:w-auto">
          <div className="relative w-full max-w-[824px] h-[300px] md:h-[400px] lg:h-[549px] bg-primary-100 rounded-lg" />
        </div>
      </div>
    </section>
  );
};
