import React, { JSX } from "react";
import { Button } from "./ui/button";

export const DataManagementSection = (): JSX.Element => {
  return (
    <section className="flex items-center justify-center px-4 md:px-8 lg:px-[220px] py-16 md:py-20 lg:py-[140px] w-full relative">
      <div className="relative w-full max-w-7xl">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12 lg:gap-16 items-center">
          <div className="flex flex-col items-start gap-8 md:gap-12 lg:gap-[60px] relative">
            <div className="flex flex-col items-start gap-4 md:gap-6 relative w-full">
              <img
                className="absolute w-[200px] md:w-[280px] lg:w-[328px] h-[29px] md:h-[38px] lg:h-[47px] top-[80px] md:top-[100px] lg:top-[129px] left-1/2 lg:left-[218px] transform -translate-x-1/2 lg:translate-x-0 z-10 hidden md:block"
                alt="Element"
                src="/element-6.png"
              />

              <h1 className="relative w-full font-header-h1 font-[number:var(--header-h1-font-weight)] text-[#212529] text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center lg:text-left tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
                Built for the Driven African Tech Professional
              </h1>

              <p className="relative w-full [font-family:'Inter',Helvetica] font-normal text-[#212529] text-base md:text-lg leading-6 md:leading-[18px] text-center lg:text-left">
                <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                  You have the degree and the foundational skills but lack the
                  "professional experience" everyone asks for. We help you frame
                  your projects and education to land that crucial first role.
                </span>
              </p>
            </div>

            <Button className="bg-[#7bbfb3] hover:bg-[#6ba89d] rounded-lg inline-flex items-center justify-center gap-2.5 px-6 py-4 md:px-10 md:py-5 h-auto mx-auto lg:mx-0">
              <span className="[font-family:'Inter',Helvetica] font-bold text-white text-base md:text-lg tracking-[0] leading-[23px]">
                Read more
              </span>
              <img className="flex-[0_0_auto]" alt="Icon" src="/icon-17.svg" />
            </Button>
          </div>

          <div className="w-full h-[300px] md:h-[400px] lg:h-[532px] bg-primary-100 rounded-lg" />
        </div>
      </div>
    </section>
  );
};
