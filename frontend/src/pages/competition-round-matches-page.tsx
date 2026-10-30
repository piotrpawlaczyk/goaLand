import { useEffect, useState } from "react";
import { CompetitionMatchCard } from "../components/competition-match-card";
import { MatchWeekHeader } from "../components/match-week-header";
import type { Edition } from "../config/editions";
import type { Round, RoundMatches } from "../config/matches";
import { useOutletContext } from "react-router";

function roundName(round: Round){
  if (round.stage == "REGULAR_SEASON"){
    return `Kolejka ${round.matchday}`
  }
  else if (round.matchday != null){
    return `${round.stage} ${round.matchday}`
  }
    return round.stage
}

type CompetitionRoundMatchesPageProps = {
  fetchMatches: (edition: Edition) => Promise<RoundMatches[]>;
};

export function CompetitionRoundMatchesPage({fetchMatches}: CompetitionRoundMatchesPageProps) {
  const [roundMatches, setRoundMatches] = useState<RoundMatches[]>([])
  const edition = useOutletContext<Edition>();
  
  useEffect(() => {
    fetchMatches(edition)
      .then(setRoundMatches)
      .catch(console.error);
  }, [edition, fetchMatches]);
  
  return(
    <div className="flex flex-col gap-[15px] pt-[50px]">{
      roundMatches.filter(({matches}) => matches.length > 0).map(({matches, round}) => (
      <CompetitionMatchCard
        key={`${round.stage}-${round.matchday ?? "all"}`}
        header={  <MatchWeekHeader  text={roundName(round)} />}
        matches={matches}
        isCurrent={edition.isCurrent}
      />
      ))
    }</div>
  );
}