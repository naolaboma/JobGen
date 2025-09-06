import { ArrowRightIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";

export const CustomizationSection = (): JSX.Element => {
  return (
    <section className="w-full flex items-center justify-center bg-[#7bbfb3] py-16 md:py-24 lg:py-[140px] px-4 md:px-8 lg:px-[220px]">
      <div className="flex flex-col lg:flex-row items-center justify-center gap-8 md:gap-12 lg:gap-[98px] max-w-full">
        <div className="flex flex-col w-full lg:w-[697px] items-start gap-8 md:gap-12 lg:gap-[60px] relative">
          <div className="flex flex-col items-start gap-4 md:gap-6 relative w-full">
            <img
              className="absolute w-[250px] md:w-[350px] lg:w-[410px] h-[50px] md:h-[65px] lg:h-[81px] top-[100px] md:top-[130px] lg:top-[164px] left-1/2 lg:left-[105px] transform -translate-x-1/2 lg:translate-x-0 z-10"
              alt="Element"
              src="/element-3.png"
            />

            <h2 className="w-full font-header-h1 font-[number:var(--header-h1-font-weight)] text-white text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center lg:text-left tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
              See Your Next Step Clearly
            </h2>

            <p className="w-full [font-family:'Inter',Helvetica] font-normal text-white text-base md:text-lg leading-6 md:leading-[18px] text-center lg:text-left">
              <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                Understand exactly what skills you need to develop. Our match
                scores and feedback highlight skill gaps, giving you a clear
                roadmap for your professional development.
              </span>
            </p>
          </div>

          <Button className="h-auto gap-2.5 px-6 py-4 md:px-10 md:py-5 bg-[#38b2ac] hover:bg-[#319795] rounded-md mx-auto lg:mx-0">
            <span className="[font-family:'Inter',Helvetica] font-bold text-white text-lg tracking-[-0.06px] leading-[23px]">
              Let&apos;s Go
            </span>
            <ArrowRightIcon className="w-4 h-4 text-white" />
          </Button>
        </div>

        <div className="flex flex-col items-start gap-1">
          <div className="w-full max-w-[500px] lg:max-w-[686px] h-[300px] md:h-[400px] lg:h-[479px] bg-primary-100 rounded-lg" />
        </div>
      </div>
    </section>
  );
};
