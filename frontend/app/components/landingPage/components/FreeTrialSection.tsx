import { ArrowRightIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";

export const FreeTrialSection = (): JSX.Element => {
  return (
    <section className="flex flex-col w-full items-center justify-center gap-12 md:gap-20 lg:gap-[200px] pt-16 md:pt-24 lg:pt-[140px] pb-8 px-4 md:px-8 lg:px-[220px] relative bg-[#7bbfb3]">
      <div className="inline-flex flex-col items-center justify-center gap-8 md:gap-10 relative flex-[0_0_auto]">
        <div className="inline-flex flex-col items-center gap-4 md:gap-6 relative flex-[0_0_auto]">
          <h1 className="relative w-full max-w-[608px] mt-[-1.00px] font-header-h1 font-[number:var(--header-h1-font-weight)] text-white text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
            Your Global Career Journey Starts Here
          </h1>

          <p className="relative w-full max-w-[550px] [font-family:'Inter',Helvetica] font-normal text-white text-lg md:text-xl lg:text-2xl text-center leading-6 md:leading-7 lg:leading-6 px-4">
            <span className="font-bold tracking-[-0.12px] leading-6 md:leading-7 lg:leading-8">
              Don&#39;t let the perfect job opportunity pass you by. Optimize
              your profile and discover your potential.
            </span>
          </p>
        </div>

        <Button className="bg-[#38b2ac] hover:bg-[#319795] rounded-[10px] inline-flex items-center justify-center gap-2.5 px-6 py-4 md:px-10 md:py-5 relative flex-[0_0_auto] h-auto">
          <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-sm md:text-base leading-4 font-bold tracking-[-0.05px]">
            Start with JobGen Today
          </span>
          <ArrowRightIcon className="relative flex-[0_0_auto] mr-[-0.50px] w-4 h-4 text-white" />
        </Button>

        <div className="relative w-[260px] h-[60px]" />
      </div>
    </section>
  );
};
