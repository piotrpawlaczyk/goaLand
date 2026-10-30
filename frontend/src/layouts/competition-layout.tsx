import { useEffect, useState } from "react";
import { Outlet, useParams } from "react-router";
import { GetCompetitionEditions } from "../api/edition";
import { CompetitionTitle } from "../components/competition-title";
import { Header } from "../components/header";
import { NavigatorButtonList } from "../components/navigator-button-list";
import type { Edition } from "../config/editions";

export function CompetitionLayout() {
  const { competitionID, startYear } = useParams();
  const [edition, setEdition] = useState<Edition>();

  useEffect(() => {
    const competitionIdNumber = Number(competitionID);
    const startYearNumber = Number(startYear);

    if (!Number.isInteger(competitionIdNumber) || !Number.isInteger(startYearNumber)) {
      setEdition(undefined);
      return;
    }

    setEdition(undefined);
    GetCompetitionEditions(competitionIdNumber)
      .then((editions) => {
        setEdition(
          editions.find((item) => item.startYear === startYearNumber),
        );
      })
      .catch(() => {
        setEdition(undefined);
      });
  }, [competitionID, startYear]);

  if (!edition) {
    return <div>Nie znaleziono rozgrywek lub trwa ładowanie</div>;
  }
  
  return (
    <main className="min-h-screen bg-background px-[50px] pb-[50px] pt-[25px]">
      <Header />
      <div className="my-[20px] flex justify-center">
        <CompetitionTitle {...edition} />
      </div>
      <div className="flex items-start">
        <NavigatorButtonList />
        <Outlet context={edition}/>
      </div>
    </main>
  );
}
