import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { get } from "../api/client";
import { Header } from "../components/header";
import { Match } from "../components/match";
import { Table } from "../components/table";
import { TeamPlayersTable } from "../components/team-players-table";
import type { ClubStats } from "../config/club";
import type { MatchData } from "../config/matches";
import type { PlayersOnPosition } from "../config/team-players";
import { ClubLogo } from "../components/club-logo";

type MatchesView = "matches" | "results";
type ClubView = "table" | "players";

type ViewButtonProps = {
  active: boolean;
  onClick: () => void;
  text: string;
};

function ViewButton({ active, onClick, text }: ViewButtonProps) {
  return (
    <button
      className={`flex h-[30px] w-[100px] items-center justify-center rounded-xs text-secondary ${
        active ? "bg-green-700" : "bg-green-750 hover:bg-green-725"
      }`}
      onClick={onClick}
      type="button"
    >
      {text}
    </button>
  );
}

export function ClubPage() {
  const { competition, clubID } = useParams<{
    competition: string;
    clubID: string;
  }>();
  const selectedClub = Number(clubID);
  const [matchesView, setMatchesView] = useState<MatchesView>("matches");
  const [clubView, setClubView] = useState<ClubView>("table");
  const [matches, setMatches] = useState<MatchData[] | null>(null);
  const [results, setResults] = useState<MatchData[] | null>(null);
  const [clubs, setClubs] = useState<ClubStats[] | null>(null);
  const [players, setPlayers] = useState<PlayersOnPosition[] | null>(null);

  
  const validParams = Boolean(competition) && Number.isInteger(selectedClub);
  const clubPath = `/competitions/${competition}/clubs/${selectedClub}`;

  useEffect(() => {
    setMatches(null);
    setResults(null);
    setClubs(null);
    setPlayers(null);
  }, [clubPath]);

  useEffect(() => {
    if (!validParams) {
      return;
    }

    const setData = matchesView === "matches" ? setMatches : setResults;
    const currentData = matchesView === "matches" ? matches : results;

    if (currentData !== null) {
      return;
    }

    get<MatchData[]>(`${clubPath}/${matchesView}`)
      .then(setData)
      .catch(console.error);
  }, [clubPath, matches, matchesView, results, validParams]);

  useEffect(() => {
    if (!validParams || clubs !== null) {
      return;
    }

    get<ClubStats[]>(`${clubPath}/table`)
      .then(setClubs)
      .catch(console.error);
  }, [clubPath, clubs, validParams]);

  useEffect(() => {
    if (!validParams || clubView !== "players" || players !== null) {
      return;
    }

    get<PlayersOnPosition[]>(`${clubPath}/players`)
      .then(setPlayers)
      .catch(console.error);
  }, [clubPath, clubView, players, validParams]);

  if (!validParams) {
    return <div>Nieprawidłowy klub</div>;
  }

  const selectedClubData = clubs?.find((club) => club.id === selectedClub);
  const clubName = selectedClubData?.name ?? `Klub ${selectedClub}`;
    
  const displayedMatches = matchesView === "matches" ? matches : results;

  return (
    <main className="min-h-screen bg-dark-background px-[50px] pb-[50px] pt-[25px]">
      <Header />

      <h2 className="mx-auto my-[40px] flex w-fit items-center gap-[10px] rounded-lg bg-card px-[50px] py-[5px] text-heading">
        <ClubLogo url={selectedClubData?.crestUrl} size={64} />
        {clubName}
      </h2>

      <div className="mx-auto flex w-fit gap-[40px]">
        <section className="flex w-[677px] shrink-0 flex-col items-center gap-[30px]">
          <nav className="flex gap-[80px]" aria-label="Mecze klubu">
            <ViewButton
              active={matchesView === "matches"}
              onClick={() => setMatchesView("matches")}
              text="mecze"
            />
            <ViewButton
              active={matchesView === "results"}
              onClick={() => setMatchesView("results")}
              text="wyniki"
            />
          </nav>

          <div className="rounded-lg bg-sections py-[10px]">
            {displayedMatches?.map((match) => (
              <Match key={match.id} {...match} size="small" />
            ))}
          </div>
        </section>

        <section className="flex w-[630px] shrink-0 flex-col items-center gap-[30px]">
          <nav className="flex gap-[80px]" aria-label="Informacje o klubie">
            <ViewButton
              active={clubView === "table"}
              onClick={() => setClubView("table")}
              text="tabela"
            />
            <ViewButton
              active={clubView === "players"}
              onClick={() => setClubView("players")}
              text="kadra"
            />
          </nav>

          {clubView === "table" && clubs && (
            <div className="rounded-lg bg-sections py-[10px]">
              <Table
                clubs={clubs}
                selectedClub={selectedClub}
                size="small"
              />
            </div>
          )}

          {clubView === "players" && players && (
            <TeamPlayersTable playersInTeam={players} />
          )}
        </section>
      </div>
    </main>
  );
}
