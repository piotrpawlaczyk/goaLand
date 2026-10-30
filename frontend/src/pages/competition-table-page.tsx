import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { GetCompetitionData } from "../api/edition";
import { getTable } from "../api/table";
import { CompetitionTitle } from "../components/competition-title";
import { Header } from "../components/header";
import { Table } from "../components/table";
import type { ClubStats } from "../config/club";
import type { Edition } from "../config/editions";

export function CompetitionTablePage() {
  const { competitionID, startYear } = useParams();
  const [edition, setEdition] = useState<Edition>();
  const [clubs, setClubs] = useState<ClubStats[]>([]);

  useEffect(() => {
    const competitionIdNumber = Number(competitionID);
    const startYearNumber = Number(startYear);

    if (!Number.isInteger(competitionIdNumber) || !Number.isInteger(startYearNumber)) {
      setEdition(undefined);
      setClubs([]);
      return;
    }

    setEdition(undefined);
    setClubs([]);

    GetCompetitionData(competitionIdNumber, startYearNumber)
      .then((selectedEdition) => {
        setEdition(selectedEdition);
        return getTable(selectedEdition);
      })
      .then(setClubs)
      .catch(console.error);
  }, [competitionID, startYear]);

  if (!edition) {
    return <div>Nie znaleziono rozgrywek</div>;
  }

  return (
    <main className="min-h-screen bg-dark-background px-[50px] pb-[50px] pt-[25px]">
      <Header />
      <div className="my-[20px] flex justify-center">
        <CompetitionTitle {...edition} />
      </div>
      <div className="flex justify-center">
        <Table clubs={clubs} isCurrent={edition.isCurrent} />
      </div>
    </main>
  );
}
