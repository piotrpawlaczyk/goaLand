import { ArrowLeft } from "lucide-react";
import { Link } from "react-router";
import { Header } from "../components/header";

export function ErrorPage() {
  return (
    <main className="flex min-h-screen flex-col bg-background px-[50px] pb-[50px] pt-[25px]">
      <Header />
      <section className="flex flex-1 items-center justify-center">
        <div className="flex flex-col items-center gap-[15px]">
          <p className="text-heading">Błąd 404</p>
          <h2 className="text-heading-bold">Nie znaleziono strony</h2>
          <Link
            className="flex items-center gap-[5px] self-start rounded-sm bg-header px-[15px] py-[5px] text-primary-bold"
            to="/"
          >
            <ArrowLeft aria-hidden="true" size={15} />
            Strona główna
          </Link>
        </div>
      </section>
    </main>
  );
}
