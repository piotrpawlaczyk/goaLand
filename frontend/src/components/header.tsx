import { Link } from "react-router";

export function Header() {
  return (
    <header className="flex h-[109px] w-full items-center justify-center gap-[50px] rounded-lg bg-header px-[77px] py-[25px]">
      <h1 className="w-[760px] text-center text-heading-bold">
        <Link
          className="inline-block rounded-sm px-[10px] py-[5px]"
          to="/"
        >
          Goaland
        </Link>
      </h1>
    </header>
  );
}
