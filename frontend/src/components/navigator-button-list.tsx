import { NavigatorButton } from "./navigator-button";

const buttonNames = ["wyniki", "mecze", "tabela", "strzelcy"] as const;

export function NavigatorButtonList() {
  return (
    <nav
      aria-label="Nawigacja rozgrywek"
      className="sticky top-[20px] flex w-[370px] flex-col items-center gap-[20px] pt-[15px] self-start"
    >
      {buttonNames.map((buttonName) => (
        <NavigatorButton
          key={buttonName}
          text={buttonName}
          to={buttonName}
        />
      ))}
    </nav>
  );
}
