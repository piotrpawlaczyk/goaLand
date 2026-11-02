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
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const competitionIdNumber = Number(competitionID);
    const startYearNumber = Number(startYear);

    if (!Number.isInteger(competitionIdNumber) || !Number.isInteger(startYearNumber)) {
      setEdition(undefined);
      setIsLoading(false);
      return;
    }

    let cancelled = false;
    setIsLoading(true);
    setEdition(undefined);
    GetCompetitionEditions(competitionIdNumber)
      .then((editions) => {
        if (!cancelled) {
          setEdition(
            editions.find((item) => item.startYear === startYearNumber),
          );
        }
      })
      .catch(() => {
        if (!cancelled) {
          setEdition(undefined);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setIsLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [competitionID, startYear]);

  if (isLoading) {
    return <div role="status">Ładowanie rozgrywek…</div>;
  }

  if (!edition) {
    return <div>Nie znaleziono rozgrywek</div>;
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
